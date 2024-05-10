// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package stor

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	. "github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/lock"
	. "github.com/donnie4w/gofer/mmap"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

var serve = &servie{}

type servie struct{}

func (t *servie) Serve() (err error) {
	return initStore()
}

func (t *servie) Close() (err error) {
	stopstat = true
	<-time.After(2 * time.Second)
	for _, ldb := range dbMap {
		ldb.Close()
	}
	dataEg.mm.Range(func(k uint64, v *Mmap) bool {
		v.Unmap()
		return true
	})
	return
}

var wfsdb *ldb
var stopstat bool
var seq int64
var count int64
var nextfn *fileHandler
var defragStat = false
var defragfile *os.File
var unmountmap = &sync.Map{}
var defragmap = &sync.Map{}

func init() {
	sys.Serve.Put(1, serve)
	sys.AppendData = fe.append
	sys.GetData = fe.getData
	sys.DelData = fe.delData
	sys.Add = fe.add
	sys.Del = fe.del
	sys.Count = fe.count
	sys.Seq = fe.seq
	sys.SearchLike = fe.findLike
	sys.SearchLimit = fe.findLimit
	sys.FragAnalysis = fe.fragAnalysis
	sys.Defrag = fe.defragAndCover
	sys.Modify = fe.modify
	sys.Import = importData
	sys.ImportFile = importFile
	sys.Export = exportData
	sys.ExportByINCR = exportByIncr
	sys.ExportByPaths = exportByPaths
	sys.ExportFile = exportFile
	sys.IsEmptyBigFile = isEmptyBigFile
}

func initStore() (err error) {
	if wfsdb, err = New(sys.WFSDATA + "/wfsdb"); err != nil {
		fmt.Println("init error:" + err.Error())
		os.Exit(1)
	}
	var wfsCurrent string
	if v, err := wfsdb.Get(CURRENT); err == nil && v != nil {
		wfsCurrent = string(v)
	}
	if v, err := wfsdb.Get(SEQ); err == nil && v != nil {
		seq = goutil.BytesToInt64(v)
	}
	if v, err := wfsdb.Get(COUNT); err == nil && v != nil {
		count = goutil.BytesToInt64(v)
	}
	initDefrag()
	if err = openFileEg(wfsCurrent); err == nil {
		initcache()
		go storTk()
	}
	return
}

func initDefrag() {
	filepath.WalkDir(sys.WFSDATA+"/wfsfile", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if strings.Contains(d.Name(), "_") {
				name := d.Name()
				name = name[:strings.Index(name, "_")]
				if id, ok := strToInt(name); ok && util.CheckNodeId(int64(id)) {
					os.Remove(path)
					fe.defragAndCover(name)
				}
			} else if id, ok := strToInt(d.Name()); ok && util.CheckNodeId(int64(id)) {
				if isEmptyBigFile(path) {
					os.Remove(path)
				}
			}
		}
		return nil
	})
}

func isEmptyBigFile(path string) bool {
	if f, err := os.Open(path); err == nil {
		defer f.Close()
		bs := make([]byte, fingerprintLen()*10)
		f.Read(bs)
		if bytes.Compare(bs, make([]byte, fingerprintLen()*10)) == 0 {
			return true
		}
	}
	return false
}

var numlock = lock.NewNumLock(1 << 9)

var dataEg = &dataHandler{mm: NewMap[uint64, *Mmap]()}

var referMap = NewLimitMap[string, *int32](1 << 15)

type dataHandler struct {
	mm *Map[uint64, *Mmap]
}

func (t *dataHandler) openMMap(node string) (_r bool) {
	if id, b := strToInt(node); b {
		lockid := goutil.CRC64(append(OPENMMAPLOCK_, goutil.Int64ToBytes(int64(id))...))
		numlock.Lock(int64(lockid))
		defer numlock.Unlock(int64(lockid))
		if !t.mm.Has(id) {
			path := getpathBynode(node)
			if goutil.IsFileExist(path) {
				if f, err := util.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666); err == nil {
					if n, err := NewMMAP(f, 0); err == nil {
						t.mm.Put(id, n)
						_r = true
					}
				}
			}
		}
		_r = t.mm.Has(id)
	}
	return
}

func (t *dataHandler) getData(node string, offset int64, size int64) (bs []byte, ok bool) {
	if defragStat {
		if _, b := defragmap.Load(node); b {
			if bs, b := t.getDataByfile(offset, size); b {
				return bs, b
			}
		}
	}
	if id, b := strToInt(node); b {
		if _, ok := unmountmap.Load(id); !ok && !t.mm.Has(id) {
			t.openMMap(node)
		}
		if m, b := t.mm.Get(id); b {
			if offset+size <= int64(len(m.Bytes())) {
				bs, ok = m.Bytes()[offset:offset+size], true
			}
		}
	}
	return
}

func (t *dataHandler) getDataByfile(offset int64, size int64) (bs []byte, ok bool) {
	defer util.Recover()
	if defragfile != nil {
		bs = make([]byte, size)
		if n, err := defragfile.ReadAt(bs, offset); err == nil && n == int(size) {
			return bs, true
		} else {
			return nil, false
		}
	}
	return
}

func (t *dataHandler) reSetMMap(node string, m *Mmap) {
	if id, b := strToInt(node); b {
		lockid := goutil.CRC64(append(RESETMMAPLOCK_, goutil.Int64ToBytes(int64(id))...))
		numlock.Lock(int64(lockid))
		defer numlock.Unlock(int64(lockid))
		if oldm, b := t.mm.Get(id); b {
			oldm.UnmapAndCloseFile()
		}
		t.mm.Put(id, m)
	}
}

func (t *dataHandler) unMmap(node string) {
	if id, b := strToInt(node); b {
		t.unMmapById(id)
	}
}

func (t *dataHandler) unMmapById(id uint64) {
	if oldm, b := t.mm.Get(id); b {
		oldm.UnmapAndCloseFile()
	}
	t.mm.Del(id)
}

var fe = &fileEg{mux: &sync.Mutex{}}

type fileEg struct {
	handler *fileHandler
	mux     *sync.Mutex
}

func openFileEg(node string) (err error) {
	fe.handler, err = initFileHandler(node)
	return
}

func (t *fileEg) add(key, value []byte) error {
	return wfsdb.Put(key, value)
}

func (t *fileEg) del(key []byte) error {
	return wfsdb.Del(key)
}

func (t *fileEg) count() int64 {
	return count
}
func (t *fileEg) seq() int64 {
	return seq
}

func (t *fileEg) findLike(pathprx string) (_r []*sys.PathBean) {
	defer util.Recover()
	pathpre := append(PATH_PRE, []byte(pathprx)...)
	if bys, err := wfsdb.GetValuesPrefix(pathpre); err == nil {
		_r = make([]*sys.PathBean, 0)
		for _, v := range bys {
			i := goutil.BytesToInt64(v)
			pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(i)...)
			if wpbbs, err := wfsdb.Get(pathseqkey); err == nil {
				wpb := bytesToWfsPathBean(wpbbs)
				if bs := t.getData(*wpb.Path); bs != nil {
					pb := &sys.PathBean{Id: i, Path: *wpb.Path, Body: bs, Timestramp: *wpb.Timestramp}
					_r = append(_r, pb)
				} else {
					t.delData(*wpb.Path)
				}
			}
		}
	}
	return
}

func (t *fileEg) findLimit(start, limit int64) (_r []*sys.PathBean) {
	if start-limit > seq {
		return
	}
	defer util.Recover()
	var count int64
	_r = make([]*sys.PathBean, 0)
	for i := start; i > 0 && count < limit; i-- {
		pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(i)...)
		if wpbbs, err := wfsdb.Get(pathseqkey); err == nil {
			wpb := bytesToWfsPathBean(wpbbs)
			if bs := t.getData(*wpb.Path); bs != nil {
				pb := &sys.PathBean{Id: i, Path: *wpb.Path, Body: bs, Timestramp: *wpb.Timestramp}
				_r = append(_r, pb)
				count++
			} else {
				t.delData(*wpb.Path)
			}
		} else if i > seq {
			count++
		}
	}
	return
}

func (t *fileEg) append(path string, bs []byte, compressType int32) (id int64, _r sys.ERROR) {
	if stopstat {
		return id, sys.ERR_STOPSERVICE
	}
	if path == "" || bs == nil || len(bs) == 0 {
		return id, sys.ERR_PARAMS
	}
	if len(bs) > int(sys.FileSize) {
		return id, sys.ERR_OVERSIZE
	}
	defer util.Recover()
	node := t.handler.Node
	var nf bool
	if nf, _r = t.handler.append(path, bs, compressType); _r != nil && _r.Equal(sys.ERR_FILEAPPEND) {
		if err := t.next(node); err == nil {
			nf, _r = t.handler.append(path, bs, compressType)
		} else {
			return id, sys.ERR_FILECREATE
		}
	}

	if nf && _r == nil && sys.Mode == 1 {
		m := make(map[*[]byte][]byte, 0)
		id = atomic.AddInt64(&seq, 1)

		m[&SEQ] = goutil.Int64ToBytes(seq)
		pathpre := append(PATH_PRE, []byte(path)...)
		m[&pathpre] = goutil.Int64ToBytes(id)

		pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(id)...)
		t := time.Now().UnixNano()
		wpb := &stub.WfsPathBean{Path: &path, Timestramp: &t}
		m[&pathseqkey] = wfsPathBeanToBytes(wpb)

		wfsdb.BatchPut(m)
	}
	return
}

func (t *fileEg) getData(path string) (_r []byte) {
	if stopstat {
		return nil
	}
	defer util.Recover()
	tasklimit()
	fidbs := fingerprint([]byte(path))
	if v, err := cacheGet(fidbs); err == nil {
		if v, err = cacheGet(v); err == nil {
			wfb := bytesToWfsFileBean(v)
			if bs, b := dataEg.getData(*wfb.Storenode, *wfb.Offset, int64(fileoffset())+*wfb.Size); b {
				_r = praseUncompress(bs[fileoffset():], *wfb.CompressType)
			}
		}
	}
	return
}

func (t *fileEg) delData(path string) (_r sys.ERROR) {
	if stopstat {
		return sys.ERR_STOPSERVICE
	}
	defer util.Recover()
	fidbs := fingerprint([]byte(path))
	batchmap := make(map[*[]byte][]byte, 0)
	dels := [][]byte{fidbs}
	if oldBidBs, err := wfsdb.Get(fidbs); err == nil && oldBidBs != nil {
		if oldWffsBs, err := wfsdb.Get(oldBidBs); err == nil && oldWffsBs != nil {
			oldwffs := bytesToWfsFileBean(oldWffsBs)
			nid, _ := strToInt(*oldwffs.Storenode)
			nidbs := goutil.Int64ToBytes(int64(nid))
			*oldwffs.Refercount -= 1
			if *oldwffs.Refercount <= 0 {
				if nodebs, err := wfsdb.Get(nidbs); err == nil && nodebs != nil {
					wnb := bytesToWfsNodeBean(nodebs)
					*wnb.Rmsize = *wnb.Rmsize + *oldwffs.Size
					batchmap[&nidbs] = wfsNodeBeanToBytes(wnb)
					dels = append(dels, oldBidBs)
				}
			} else {
				batchmap[&oldBidBs] = wfsFileBeanToBytes(oldwffs)
			}
		}
		if sys.Mode == 1 {
			pathpre := append(PATH_PRE, []byte(path)...)
			if v, err := wfsdb.Get(pathpre); err == nil {
				dels = append(dels, pathpre)
				dels = append(dels, append(PATH_SEQ, v...))
			}
		}
		batchmap[&COUNT] = goutil.Int64ToBytes(atomic.AddInt64(&count, -1))
	} else {
		return sys.ERR_NOTEXSIT
	}

	if err := wfsdb.Batch(batchmap, dels); err != nil {
		return sys.ERR_UNDEFINED
	} else {
		cacheDel(fidbs)
	}
	return
}

func (t *fileEg) modify(path, newpath string) (err sys.ERROR) {
	if stopstat {
		return sys.ERR_STOPSERVICE
	}
	if path == "" || newpath == "" || path == newpath {
		return sys.ERR_PARAMS
	}
	am := make(map[*[]byte][]byte, 0)
	dm := make([][]byte, 0)

	fidbs := fingerprint([]byte(path))
	dm = append(dm, fidbs)

	newfidbs := fingerprint([]byte(newpath))

	if sys.Mode == 1 {
		pathpre := append(PATH_PRE, []byte(path)...)
		dm = append(dm, pathpre)
		if v, err := wfsdb.Get(pathpre); err == nil {
			newpathpre := append(PATH_PRE, []byte(newpath)...)
			am[&newpathpre] = v

			i := goutil.BytesToInt64(v)
			pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(i)...)
			if v, err := wfsdb.Get(pathseqkey); err == nil {
				wpb := bytesToWfsPathBean(v)
				wpb.Path = &newpath
				am[&pathseqkey] = wfsPathBeanToBytes(wpb)
			}
			cacheDel(pathpre)
		} else {
			return sys.ERR_NOTEXSIT
		}
	}

	if oldBidBs, err := wfsdb.Get(fidbs); err == nil && oldBidBs != nil {
		am[&newfidbs] = oldBidBs
	} else {
		return sys.ERR_NOTEXSIT
	}
	if _, err := wfsdb.Get(newfidbs); err != nil && len(am) > 0 && len(dm) > 0 {
		if wfsdb.Batch(am, dm) == nil {
			cacheDel(fidbs)

		}
	} else {
		return sys.ERR_NEWPATHEXIST
	}
	return
}

var atomicflag int32 = 0

func (t *fileEg) next(node string) (err error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	if node == t.handler.Node {
		if nextfn != nil {
			t.handler = nextfn
			usefileHandler(t.handler)
			nextfn = nil
		} else {
			t.handler, err = initFileHandler("")
		}
	}
	return
}

func newNextfn() {
	if atomic.CompareAndSwapInt32(&atomicflag, 0, 1) {
		defer atomic.SwapInt32(&atomicflag, 0)
		if nextfn == nil {
			nextfn, _ = newFileHandler()
		}
	}
}

// Deprecated
func (t *fileEg) defrag(node string) (err sys.ERROR) {
	if stopstat {
		return sys.ERR_STOPSERVICE
	}
	defragStat = true
	defer func() {
		if e := recover(); e != nil {
			err = sys.ERR_UNDEFINED
		}
		defragStat = false
		defragmap.Delete(node)
	}()
	if v, err := wfsdb.Get(CURRENT); err == nil && v != nil {
		if string(v) == node {
			return sys.ERR_DEFRAG_FORBID
		}
	}

	nid, b := strToInt(node)
	if !b {
		return sys.ERR_NOTEXSIT
	}
	dataEg.openMMap(node)
	mm, b := dataEg.mm.Get(nid)
	if !b {
		return sys.ERR_NOTEXSIT
	}

	newnid := util.CreateNodeId()
	newnode := intToStr(uint64(newnid))
	newnidbs := goutil.Int64ToBytes(newnid)
	newnodepath := getpathBynode(newnode)
	if err := os.MkdirAll(filepath.Dir(newnodepath), 0777); err != nil {
		return sys.ERR_UNDEFINED
	}
	if f, err := util.OpenFile(newnodepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err == nil {
		defragfile = f
		defer f.Close()
		wl := new(int64)
		if err = defragFB(mm.Bytes(), f, 0, &newnode, wl, new(int32)); err == nil {
			var size int64
			if fs, _ := f.Stat(); fs != nil {
				size = fs.Size()
			} else {
				size = *wl
			}
			oldOffsBs := append(ENDOFFSET_, goutil.Int64ToBytes(int64(nid))...)
			newOffset := append(ENDOFFSET_, goutil.Int64ToBytes(int64(newnid))...)
			newbat := make(map[*[]byte][]byte, 0)
			if *wl == 0 && size == 0 {
				f.Close()
				os.Remove(newnodepath)
				f = nil
			} else {
				newbat[&newOffset] = goutil.Int64ToBytes(size)
				newbat[&newnidbs] = wfsNodeBeanToBytes(&stub.WfsNodeBean{Rmsize: new(int64)})
			}
			wfsdb.Batch(newbat, [][]byte{oldOffsBs})
			dataEg.unMmap(node)
			if err := os.Remove(getpathBynode(node)); err != nil {
				logging.Error(err)
			}
			dataEg.unMmap(newnode)
			dataEg.openMMap(newnode)

		} else {
			if fs, e := f.Stat(); e == nil {
				if fs.Size() == 0 && *wl == 0 {
					f.Close()
					os.Remove(newnodepath)
					f = nil
				}
			}
			return sys.ERR_UNDEFINED
		}
	}
	return
}

func (t *fileEg) defragAndCover(node string) (err sys.ERROR) {
	if stopstat {
		return sys.ERR_STOPSERVICE
	}
	defragStat = true
	defer func() {
		if e := recover(); e != nil {
			err = sys.ERR_UNDEFINED
		}
		defragStat = false
		defragmap.Delete(node)
	}()
	if v, err := wfsdb.Get(CURRENT); err == nil && v != nil {
		if string(v) == node {
			return sys.ERR_DEFRAG_FORBID
		}
	}
	nid, b := strToInt(node)
	if !b {
		return sys.ERR_NOTEXSIT
	}
	dataEg.openMMap(node)
	mm, b := dataEg.mm.Get(nid)
	if !b {
		return sys.ERR_NOTEXSIT
	}
	nidbs := goutil.Int64ToBytes(int64(nid))
	nodepath := getpathBynode(node)

	newnode := fmt.Sprint(node, "_", util.CreateNodeId())
	newnodepath := getpathBynode(newnode)
	if err := os.MkdirAll(filepath.Dir(newnodepath), 0777); err != nil {
		return sys.ERR_UNDEFINED
	}
	if f, err := util.OpenFile(newnodepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666); err == nil {
		defragfile = f
		defer f.Close()
		wl := new(int64)
		defragmap.Store(node, newnode)
		if err = defragFB(mm.Bytes(), f, 0, &node, wl, new(int32)); err == nil {
			var size int64
			if fs, _ := f.Stat(); fs != nil {
				size = fs.Size()
			} else {
				size = *wl
			}
			if *wl == 0 && size == 0 {
				f.Close()
				os.Remove(newnodepath)
				f = nil
				return sys.ERR_UNDEFINED
			} else {
				offsBs := append(ENDOFFSET_, goutil.Int64ToBytes(int64(nid))...)
				newbat := make(map[*[]byte][]byte)
				newbat[&offsBs] = goutil.Int64ToBytes(size)
				newbat[&nidbs] = wfsNodeBeanToBytes(&stub.WfsNodeBean{Rmsize: new(int64)})
				wfsdb.BatchPut(newbat)
			}
			unmountmap.Store(nid, byte(0))
			defer unmountmap.Delete(nid)
			dataEg.unMmap(node)
			dataEg.unMmap(newnode)
			f.Close()
			if os.Rename(newnodepath, nodepath) == nil {
				dataEg.openMMap(node)
			}
		} else {
			if fs, e := f.Stat(); e == nil {
				if fs.Size() == 0 && *wl == 0 {
					f.Close()
					os.Remove(newnodepath)
					f = nil
				}
			}
			return sys.ERR_UNDEFINED
		}
	}
	return
}

func getpathBynode(node string) string {
	return sys.WFSDATA + "/wfsfile/" + node
}

func defragFB(bs []byte, f *os.File, offset int64, node *string, wl *int64, rl *int32) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	step := fingerprintLen()
	if len(bs) > step+4 {
		if v, e := wfsdb.Get(bs[:step]); e == nil && v != nil {
			wfb := bytesToWfsFileBean(v)
			wfb.Storenode = node
			*wfb.Offset = offset
			size := goutil.BytesToInt32(bs[step : step+4])
			if n, e := f.Write(bs[:step+4+int(size)]); e == nil {
				atomic.AddInt64(wl, int64(n))
				cacheDel(bs[:step])
				if err = wfsdb.Put(bs[:step], wfsFileBeanToBytes(wfb)); err == nil {
					return defragFB(bs[n:], f, offset+int64(n), node, wl, rl)
				}
			} else {
				return e
			}
		} else {
			if size := goutil.BytesToInt32(bs[step : step+4]); size > 0 || *rl < 4 {
				if size == 0 {
					atomic.AddInt32(rl, 1)
				}
				err = defragFB(bs[step+4+int(size):], f, offset, node, wl, rl)
			}
		}
	}
	return
}

func (t *fileEg) fragAnalysis(node string) (fb *sys.FragBean, err sys.ERROR) {
	if stopstat {
		return nil, sys.ERR_STOPSERVICE
	}
	defer func() {
		if e := recover(); e != nil {
			err = sys.ERR_UNDEFINED
		}
	}()

	if v, err := wfsdb.Get(CURRENT); err == nil && v != nil {
		if string(v) == node {
			return nil, sys.ERR_DEFRAG_FORBID
		}
	}
	fb = &sys.FragBean{Node: node}
	nid, _ := strToInt(node)
	nidbs := goutil.Int64ToBytes(int64(nid))
	ofsBs := append(ENDOFFSET_, nidbs...)

	if v, err := wfsdb.Get(ofsBs); err == nil {
		fb.ActualSize = goutil.BytesToInt64(v)
	}
	if nodebs, err := wfsdb.Get(nidbs); err == nil && nodebs != nil {
		wnb := bytesToWfsNodeBean(nodebs)
		fb.RmSize = *wnb.Rmsize
	}
	if f, err := os.Stat(getpathBynode(node)); err == nil {
		fb.FileSize = f.Size()
	}
	return
}

type fileHandler struct {
	mm     *Mmap
	Node   string
	length int64
}

func initFileHandler(node string) (fh *fileHandler, err error) {
	fine := false
	if node == "" {
		if fh, err = newFileHandler(); err == nil {
			err = usefileHandler(fh)
			return
		} else {
			sys.FmtLog(err)
		}
	} else {
		nid, _ := strToInt(node)
		nidbs := goutil.Int64ToBytes(int64(nid))
		if endoffsetBs, err := wfsdb.Get(append(ENDOFFSET_, nidbs...)); err == nil {
			if endOffset := goutil.BytesToInt64(endoffsetBs); endOffset < sys.FileSize {
				nodepath := getpathBynode(node)
				if goutil.IsFileExist(nodepath) {
					if f, err := util.OpenFile(nodepath, os.O_CREATE|os.O_RDWR, 0666); err == nil {
						if mm, err := NewMMAP(f, endOffset); err == nil {
							dataEg.reSetMMap(node, mm)
							fh = &fileHandler{mm: mm, Node: node, length: endOffset}
							fine = true
						}
					}
				}
			}
		}

		if !fine {
			initFileHandler("")
			return
		}
	}
	if !fine {
		return initFileHandler("")
	}
	return
}

func newFileHandler() (fh *fileHandler, err error) {
	nid := util.CreateNodeId()
	node := intToStr(uint64(nid))
	var f *os.File
	nodepath := getpathBynode(node)
	if err = os.MkdirAll(filepath.Dir(nodepath), 0777); err != nil {
		return
	}
	if f, err = util.OpenFile(nodepath, os.O_CREATE|os.O_RDWR, 0666); err == nil {
		if err = f.Truncate(sys.FileSize); err == nil {
			var mm *Mmap
			if mm, err = NewMMAP(f, 0); err == nil {
				dataEg.reSetMMap(node, mm)
				fh = &fileHandler{mm: mm, Node: node}
			}
		}
	}
	return
}

func usefileHandler(fh *fileHandler) (err error) {
	if nid, b := strToInt(fh.Node); b {
		nidbs := goutil.Int64ToBytes(int64(nid))
		fmap := make(map[*[]byte][]byte, 0)
		ofsBs := append(ENDOFFSET_, nidbs...)
		fmap[&ofsBs] = []byte{0}
		fmap[&nidbs] = wfsNodeBeanToBytes(&stub.WfsNodeBean{Rmsize: new(int64)})
		fmap[&CURRENT] = []byte(fh.Node)
		err = wfsdb.BatchPut(fmap)
	} else {
		return errors.New("format error")
	}
	return
}

func (t *fileHandler) append(path string, bs []byte, compressType int32) (nf bool, _r sys.ERROR) {
	if path != "" && bs != nil && len(bs) > 0 {
		fidBs := fingerprint([]byte(path))

		lockid := goutil.CRC64(append(APPENDLOCK_, fidBs...))
		numlock.Lock(int64(lockid))
		defer numlock.Unlock(int64(lockid))

		bidBs := fingerprint(bs)

		nid, _ := strToInt(t.Node)
		nidbs := goutil.Int64ToBytes(int64(nid))

		var wfbbs []byte

		if v, err := wfsdb.Get(bidBs); err != nil || v == nil {
			storeBytes := praseCompress(bs, compressType)
			if cl := atomic.AddInt64(&t.length, int64(len(storeBytes)+fileoffset())); cl < sys.FileSize {

				//when the ratio(90%) is exceeded, an empty big file will be created to avoid lock contention
				//that occurs when files are created at high concurrency.
				if nextfn == nil && float32(cl)/float32(sys.FileSize) > 0.9 {
					go newNextfn()
				}

				size, refer := int64(len(storeBytes)), new(int32)
				*refer = 1

				if r, ok := referMap.LoadOrStore(string(bidBs), refer); ok {
					refer = r
					atomic.AddInt32(refer, 1)
				}

				wfb := &stub.WfsFileBean{Storenode: &t.Node, Size: &size, CompressType: &compressType, Refercount: refer}

				fmap := make(map[*[]byte][]byte, 0)

				bs := append(bidBs, goutil.Int32ToBytes(int32(size))...)
				if !sys.SYNC {
					if n, err := t.mm.Append(append(bs, storeBytes...)); err == nil {
						wfb.Offset = &n
					} else {
						return nf, sys.ERR_FILEAPPEND
					}
				} else {
					if n, err := t.mm.AppendSync(append(bs, storeBytes...)); err == nil {
						wfb.Offset = &n
					} else {
						return nf, sys.ERR_FILEAPPEND
					}
				}

				wfbbytes := wfsFileBeanToBytes(wfb)
				fmap[&bidBs] = wfbbytes

				ofsBs := append(ENDOFFSET_, nidbs...)
				fmap[&ofsBs] = goutil.Int64ToBytes(t.length)

				if err := wfsdb.BatchPut(fmap); err != nil {
					return nf, sys.ERR_UNDEFINED
				} else {
					cachePut(bidBs, wfbbytes)
				}

			} else {
				return nf, sys.ERR_FILEAPPEND
			}
		} else {
			wfbbs = v
		}
		batchmap := make(map[*[]byte][]byte, 0)
		var dels [][]byte
		if oldBidBs, err := wfsdb.Get(fidBs); err == nil && oldBidBs != nil {
			if bytes.Equal(oldBidBs, bidBs) {
				return nf, sys.ERR_EXSIT
			}
			if oldWffsBs, err := wfsdb.Get(oldBidBs); err == nil && oldWffsBs != nil {
				oldwffs := bytesToWfsFileBean(oldWffsBs)
				*oldwffs.Refercount -= 1
				if *oldwffs.Refercount <= 0 {
					if nodebs, err := wfsdb.Get(nidbs); err == nil && nodebs != nil {
						wnb := bytesToWfsNodeBean(nodebs)
						*wnb.Rmsize = *wnb.Rmsize + *oldwffs.Size
						batchmap[&nidbs] = wfsNodeBeanToBytes(wnb)
						dels = [][]byte{oldBidBs}
					}
				} else {
					batchmap[&oldBidBs] = wfsFileBeanToBytes(oldwffs)
				}
			}
		} else {
			nf = true
			batchmap[&COUNT] = goutil.Int64ToBytes(atomic.AddInt64(&count, 1))
		}

		if wfbbs != nil {
			wfb := bytesToWfsFileBean(wfbbs)
			if r, ok := referMap.LoadOrStore(string(bidBs), wfb.Refercount); ok {
				wfb.Refercount = r
			}
			atomic.AddInt32(wfb.Refercount, 1)
			batchmap[&bidBs] = wfsFileBeanToBytes(wfb)
		}

		batchmap[&fidBs] = bidBs
		if err := wfsdb.Batch(batchmap, dels); err != nil {
			return nf, sys.ERR_UNDEFINED
		} else {
			cachePut(fidBs, bidBs)
		}
	} else {
		return nf, sys.ERR_PARAMS
	}
	return
}

func exportData(streamfunc func(bean *stub.SnapshotBean) bool) (err error) {
	defer util.Recover()
	return wfsdb.SnapshotToStream(nil, streamfunc)
}

func exportByIncr(start, limit int64, streamfunc func(snaps *stub.SnapshotBeans) bool) (err sys.ERROR) {
	defer util.Recover()
	if start > 0 && limit > 0 && start <= seq {
		nodemap := make(map[string]string)
		count := 0
		for i := start; i <= seq && count < int(limit); i++ {
			seqbs := append(PATH_SEQ, goutil.Int64ToBytes(i)...)
			if wpbbs, err := wfsdb.Get(seqbs); err == nil {
				count++
				snaps := &stub.SnapshotBeans{Id: new(int64)}
				*snaps.Id = i
				snaps.Beans = make([]*stub.SnapshotBean, 0)
				snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: seqbs, Value: wpbbs})
				wpbtb := bytesToWfsPathBean(wpbbs)
				path := *wpbtb.Path
				fidbs := fingerprint([]byte(path))
				if bidBs, err := wfsdb.Get(fidbs); err == nil {
					snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: fidbs, Value: bidBs})
					if wfbbs, err := wfsdb.Get(bidBs); err == nil {
						snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: bidBs, Value: wfbbs})
						wfb := bytesToWfsFileBean(wfbbs)
						node := *wfb.Storenode
						if _, ok := nodemap[node]; !ok {
							nodemap[node] = ""
							nid, _ := strToInt(node)
							nidbs := goutil.Int64ToBytes(int64(nid))
							if wnbbs, err := wfsdb.Get(nidbs); err == nil {
								snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: nidbs, Value: wnbbs})
							}
							ofsBs := append(ENDOFFSET_, nidbs...)
							if v, err := wfsdb.Get(ofsBs); err == nil {
								snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: ofsBs, Value: v})
							}
						}
					}
					streamfunc(snaps)
				}
			}
		}
	} else {
		return sys.ERR_PARAMS
	}
	return
}

func exportByPaths(paths []string, streamfunc func(snaps *stub.SnapshotBeans) bool) (err sys.ERROR) {
	defer util.Recover()
	if len(paths) > 0 {
		nodemap := make(map[string]string)
		for _, path := range paths {
			fidbs := fingerprint([]byte(path))
			if bidBs, err := wfsdb.Get(fidbs); err == nil {
				snaps := &stub.SnapshotBeans{Id: new(int64)}
				snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: fidbs, Value: bidBs})
				if wfbbs, err := wfsdb.Get(bidBs); err == nil {
					snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: bidBs, Value: wfbbs})
					wfb := bytesToWfsFileBean(wfbbs)
					node := *wfb.Storenode
					if _, ok := nodemap[node]; !ok {
						nodemap[node] = ""
						nid, _ := strToInt(node)
						nidbs := goutil.Int64ToBytes(int64(nid))
						if wnbbs, err := wfsdb.Get(nidbs); err == nil {
							snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: nidbs, Value: wnbbs})
						}
						ofsBs := append(ENDOFFSET_, nidbs...)
						if v, err := wfsdb.Get(ofsBs); err == nil {
							snaps.Beans = append(snaps.Beans, &stub.SnapshotBean{Key: ofsBs, Value: v})
						}
					}
				}
				streamfunc(snaps)
			}
		}
	} else {
		return sys.ERR_PARAMS
	}
	return
}

func exportFile(start, limit int64, streamfunc func(snaps *stub.SnapshotFile) bool) (err sys.ERROR) {
	if start > 0 && limit > 0 && start <= seq {
		count := int64(0)
		for i := start; i <= seq && count < limit; i++ {
			pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(i)...)
			if wpbbs, err := wfsdb.Get(pathseqkey); err == nil {
				wpb := bytesToWfsPathBean(wpbbs)
				if bs := fe.getData(*wpb.Path); bs != nil {
					count++
					fidbs := fingerprint([]byte(*wpb.Path))
					compressType := new(int32)
					if bidBs, err := wfsdb.Get(fidbs); err == nil {
						if wfbtb, err := wfsdb.Get(bidBs); err == nil {
							if wffb := bytesToWfsFileBean(wfbtb); wffb != nil {
								compressType = wffb.CompressType
							}
						}
					}
					sf := &stub.SnapshotFile{Id: &i, Path: wpb.Path, Data: bs, CompressType: compressType}
					streamfunc(sf)
				}
			}
		}
	} else {
		return sys.ERR_PARAMS
	}
	return
}

func importData(bean *stub.SnapshotBean, cover bool) (err error) {
	defer util.Recover()
	if bean == nil {
		return errors.New("data is nil")
	}
	if bytes.Equal(bean.Key, CURRENT) || bytes.Equal(bean.Key, COUNT) || bytes.Equal(bean.Key, SEQ) {
		return
	}

	if len(bean.Key) == 17 && bytes.Equal(bean.Key[:9], PATH_SEQ) {
		wppb := bytesToWfsPathBean(bean.Value)
		path := *wppb.Path
		pathpre := append(PATH_PRE, []byte(path)...)
		am := make(map[*[]byte][]byte, 0)
		dm := make([][]byte, 0)
		if v, e := wfsdb.Get(pathpre); e == nil {
			if !cover {
				return
			}
			oldpathseqkey := append(PATH_SEQ, v...)
			dm = append(dm, oldpathseqkey)
		} else {
			am[&COUNT] = goutil.Int64ToBytes(atomic.AddInt64(&count, 1))
		}
		id := atomic.AddInt64(&seq, 1)
		am[&SEQ] = goutil.Int64ToBytes(seq)
		am[&pathpre] = goutil.Int64ToBytes(id)
		pathseqkey := append(PATH_SEQ, goutil.Int64ToBytes(id)...)
		am[&pathseqkey] = bean.Value
		wfsdb.Batch(am, dm)
		return
	}

	if len(bean.Key) > 2 && bytes.Equal(bean.Key[:2], PATH_PRE) && len(bean.Value) == 8 {
		paths := bean.Key[2:]
		b := true
		for _, v := range paths {
			if v > unicode.MaxASCII {
				b = false
			}
		}
		if b && goutil.BytesToInt64(bean.Value) < 1<<50 {
			return
		}
	}
	if len(bean.Key) == fingerprintLen() {
		cacheDel(bean.Key)
	}
	err = wfsdb.LoadSnapshotBean(bean)
	return
}

func importFile(snapsBean *stub.SnapshotFile) (err sys.ERROR) {
	defer util.Recover()
	if snapsBean.Path != nil && *snapsBean.Path != "" && len(snapsBean.Data) > 0 {
		_, err = fe.append(snapsBean.GetPath(), snapsBean.GetData(), snapsBean.GetCompressType())
	}
	return
}
