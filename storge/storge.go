/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package storge

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/donnie4w/wfs/conf"
	. "github.com/donnie4w/wfs/db"
	"github.com/golang/groupcache/lru"
)

var _del_ = "_del_"
var _current_file_ = "_current_file_"
var _file_sequence_ = "_file_sequence_"
var _ERR_CODE_APPEND_DATA = 501

var db *DB
var fm *FileManager

func Init() {
	db = NewDB(CF.Data, true)
	fm = OpenFileManager()
	go _Ticker(1800, _compact)
}

func AppendData(bs []byte, name string, fileType string) (err error) {
	defer catchError()
	if name == "" || bs == nil || len(bs) == 0 {
		return errors.New("nil")
	}
	fingerprint := _fingerprint([]byte(name))
	lockString.Lock(fingerprint)
	defer lockString.UnLock(fingerprint)
	md5key := MD5(bs)
	sbs := NewSegment(name, fileType, md5key)
	DBPutSegment([]byte(fingerprint), *sbs)
	fm.setNameCache(fingerprint, sbs)
	lockString.Lock(string(md5key))
	defer lockString.UnLock(string(md5key))
	mb := fm.GetMd5Bean(md5key)
	if mb == nil {
		f := fm.getFData()
		//		sequence := fm.getSequence()
		//		bss := make([]byte, len(bs)+len(sequence))
		//		copy(bss, bs)
		//		copy(bss[len(bs):], sequence)
		//		offset := f.GetAndSetCurPoint(int64(len(bss)))
		//		mb = NewMd5Bean(offset, int32(len(bss)), f.FileName, sequence)
		//		err = f.AppendData(bss, offset)
		//		err = f.WriteIdxMd5(md5key)
		//		err = DBPut(sequence, md5key)
		offset := f.GetAndSetCurPoint(int64(len(bs)))
		mb = NewMd5Bean(offset, int32(len(bs)), f.FileName, nil)
		err = f.AppendData(bs, offset)
		err = f.WriteIdxMd5(md5key)
	}
	mb.AddQuote()
	//	fmt.Println("quote===>", mb.QuoteNum)
	DBPutMd5Bean(md5key, *mb)
	return
}

func _AppendData(bs []byte, f *Fdata) (err error) {
	//	fmt.Println("_AppendData:", f)
	md5key := MD5(bs)
	offset := f.GetAndSetCurPoint(int64(len(bs)))
	mb := NewMd5Bean(offset, int32(len(bs)), f.FileName, nil)
	err = f.AppendData(bs, offset)
	err = f.WriteIdxMd5(md5key)
	mb.AddQuote()
	DBPutMd5Bean(md5key, *mb)
	return
}

func GetData(name string) (bs []byte, er error) {
	defer catchError()
	if name == "" {
		return nil, errors.New("nil")
	}
	fingerprint := _fingerprint([]byte(name))
	segment, err := fm.getSegment(fingerprint)
	if err == nil {
		md5key := segment.Md5
		md5Bean, err := DBGetMd5Bean(md5key)
		if err == nil {
			filename := md5Bean.FileName
			//			fmt.Println("GetData:", filename)
			fdata := fm.GetFdataByName(filename)
			bs, er = fdata.GetData(&md5Bean)
			//			bs = bs[:len(bs)-8]
		}
	} else {
		er = err
	}
	return
}

func DelData(name string) (err error) {
	defer catchError()
	if name == "" {
		return errors.New("nil")
	}
	fingerprint := _fingerprint([]byte(name))
	lockString.Lock(fingerprint)
	defer lockString.UnLock(fingerprint)
	segment, er := fm.getSegment(fingerprint)
	fm.removeNameCache(fingerprint)
	err = DBDel([]byte(fingerprint))
	if er == nil {
		md5key := segment.Md5
		lockString.Lock(string(md5key))
		defer lockString.UnLock(string(md5key))
		mb := fm.GetMd5Bean(md5key)
		mb.SubQuote()
		//		fmt.Println("===>", mb.QuoteNum)
		if mb.QuoteNum <= 0 {
			fm.DelMd5Bean(md5key)
			//			fmt.Println("deldata===>", mb.Sequence)
			_saveDel(mb)
		} else {
			DBPutMd5Bean(md5key, *mb)
		}
	}
	return
}

func Exsit(name string) (b bool) {
	defer catchError()
	if name == "" {
		return
	}
	fingerprint := _fingerprint([]byte(name))
	return fm.hasName(fingerprint)
}

func _saveDel(mb *Md5Bean) {
	filename := mb.FileName
	lockString.Lock(filename)
	defer lockString.UnLock(filename)
	size := mb.Size
	filekey := fmt.Sprint(_del_, filename)
	v, err := db.Get([]byte(filekey))
	if err == nil {
		i := Bytes2Octal(v) + size
		db.Put([]byte(filekey), Octal2bytes(i))
	} else {
		db.Put([]byte(filekey), Octal2bytes(size))
	}
}

//-----------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------

var lockString = &LockString{locks: make(map[string]*sync.Mutex, 0), lockCount: make(map[string]int32, 0), lock: new(sync.Mutex)}

type LockString struct {
	locks     map[string]*sync.Mutex
	lockCount map[string]int32
	lock      *sync.Mutex
}

func (this *LockString) Lock(s string) {
	this.lock.Lock()
	var lock *sync.Mutex
	var ok bool
	if lock, ok = this.locks[s]; !ok {
		lock = new(sync.Mutex)
		this.locks[s] = lock
		this.lockCount[s] = 1
	} else {
		this.lockCount[s] = this.lockCount[s] + 1
	}
	//	fmt.Println("LockString==>", this.lockCount[s])
	this.lock.Unlock()
	lock.Lock()
}

func (this *LockString) UnLock(s string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if lock, ok := this.locks[s]; ok {
		if this.lockCount[s] == 1 {
			//			fmt.Println("UnLock lockString==>", s)
			delete(this.locks, s)
			delete(this.lockCount, s)
		} else {
			this.lockCount[s] = this.lockCount[s] - 1
		}
		lock.Unlock()
	}
}

//-------------------------------------------------------------------------------------------------------------------

type FileManager struct {
	lock         *sync.RWMutex
	fileMap      map[string]*Fdata
	fileMaxSize  int64  //
	fileSequence []byte //
	md5map       map[string]*Md5Bean
	nameCache    *lru.Cache
	currFileName string
}

func OpenFileManager() (f *FileManager) {
	f = &FileManager{lock: new(sync.RWMutex), fileMap: make(map[string]*Fdata, 0), fileMaxSize: CF.MaxDataSize, md5map: make(map[string]*Md5Bean, 0), nameCache: lru.New(1 << 20)}
	f.getFData()
	return
}

func (this *FileManager) getSegment(fingerprint string) (segment Segment, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	s, b := this.nameCache.Get(fingerprint)
	if b && s != nil {
		s1 := s.(*Segment)
		segment = *s1
		return
	}
	segment, err = DBGetSegment([]byte(fingerprint))
	if err == nil {
		this.nameCache.Add(fingerprint, &segment)
	}
	return
}
func (this *FileManager) hasName(name string) (b bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	_, b = this.nameCache.Get(name)
	if b {
		return
	}
	b = DBExsit([]byte(name))
	if b {
		s, err := DBGetSegment([]byte(name))
		if err == nil {
			this.nameCache.Add(name, &s)
		}
	}
	return
}

func (this *FileManager) setNameCache(fingerprint string, segment *Segment) {
	this.nameCache.Add(fingerprint, segment)
}

func (this *FileManager) removeNameCache(fingerprint string) {
	this.nameCache.Remove(fingerprint)
}

func (this *FileManager) getFData() (fdata *Fdata) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if filedata, ok := this.fileMap[_current_file_]; ok {
		if filedata.FileSize() < this.fileMaxSize {
			return filedata
		}
	} else {
		v, err := db.Get([]byte(_current_file_))
		var filename string
		if err == nil && v != nil {
			filename = string(v)
			fdata = this._openFdataFile(filename)
			if fdata.FileSize() < this.fileMaxSize {
				this.fileMap[filename] = fdata
				this.fileMap[_current_file_] = fdata
				db.Put([]byte(_current_file_), []byte(filename))
				this.currFileName = filename
				return
			}
		}
	}
	return this._newFdata(true)
}

//func (this *FileManager) getSequence() []byte {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	var err error
//	if this.fileSequence == nil {
//		this.fileSequence, err = db.Get([]byte(_file_sequence_))
//		if err != nil || this.fileSequence == nil {
//			this.fileSequence = []byte{0}
//		}
//	}
//	this.fileSequence = Hex2bytes(Bytes2hex(this.fileSequence) + 1)
//	db.Put([]byte(_file_sequence_), this.fileSequence)
//	return this.fileSequence
//}

func (this *FileManager) _newFdata(isCurrent bool) (fdata *Fdata) {
	filename := fmt.Sprint("data/", time.Now().Unix(), ".dat")
	fdata = this._openFdataFile(filename)
	this.fileMap[fdata.FileName] = fdata
	if isCurrent {
		this.fileMap[_current_file_] = fdata
		db.Put([]byte(_current_file_), []byte(fdata.FileName))
		this.currFileName = filename
	}
	return
}

func (this *FileManager) GetFdataByName(filename string) (fdata *Fdata) {
	this.lock.Lock()
	defer this.lock.Unlock()
	var ok bool
	if fdata, ok = this.fileMap[filename]; !ok {
		fdata = this._openFdataFile(filename)
		this.fileMap[filename] = fdata
	}
	return
}

func (this *FileManager) _openFdataFile(filename string) (fdata *Fdata) {
	//	fmt.Println("_newFdataFile=>", filename)
	idxfilename := strings.Replace(filename, ".dat", ".idx", -1)
	currFile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	currIdxFile, err := os.OpenFile(idxfilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err == nil {
		stat, _ := currFile.Stat()
		//		fmt.Println("stat.Size()===>", stat.Size())
		fdata = &Fdata{filename, stat.Size(), currFile, currIdxFile, new(sync.RWMutex)}
	} else {
		fmt.Println("errer==>", err.Error())
	}
	return
}

func (this *FileManager) GetMd5Bean(md5key []byte) (mb *Md5Bean) {
	var ok bool
	if mb, ok = this.md5map[string(md5key)]; ok {
		return
	}
	nmb, err := DBGetMd5Bean(md5key)
	if err == nil {
		mb = &nmb
		this.md5map[string(md5key)] = mb
	}
	return
}

//注意这个删除并非线程安全，mk5key并发调用仍然存储相同md5值刚存储就被删除的情况
//由于正常情况下出现机率极低，此处不做处理
func (this *FileManager) DelMd5Bean(md5key []byte) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.md5map, string(md5key))
	DBDel(md5key)
	return
}

//-------------------------------------------------------------------------------------------------------------------

type Fdata struct {
	FileName string   //所在文件名
	CurPoint int64    //当前指针
	f        *os.File //
	idxf     *os.File //
	lock     *sync.RWMutex
}

//func (this *Fdata) GetCurrSegmentId() (Id int64) {
//	this.lock.Lock()
//	defer this.lock.Unlock()
//	v, err := db.Get([]byte(fmt.Sprint("_idx_id_", this.FileName)))
//	if err == nil && v != nil {
//		Id = Bytes2hex(v) + 1
//	} else {
//		Id = 1
//	}
//	db.Put([]byte(fmt.Sprint("_idx_id_", this.FileName)), Hex2bytes(Id))
//	return
//}

func (this *Fdata) GetAndSetCurPoint(size int64) (offset int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	offset = this.CurPoint
	this.CurPoint = offset + size
	return
}

func (this *Fdata) FileSize() int64 {
	return this.CurPoint
}

func (this *Fdata) CloseFile() {
	this.f.Close()
}

func (this *Fdata) AppendData(bs []byte, offset int64) (err error) {
	//	fmt.Println("AppendData==>", this.f.Name(), " ,", len(bs), " ,", offset)
	_, err = Append(this.f, bs, offset)
	if err != nil {
		panic(_ERR_CODE_APPEND_DATA)
	}
	return
}

func (this *Fdata) WriteIdxMd5(md5key []byte) (err error) {
	_, err = Write(this.idxf, md5key)
	return
}

func (this *Fdata) CloseAndDelete() {
	fmt.Println("CloseAndDelete:", this)
	this.lock.Lock()
	defer this.lock.Unlock()
	defer catchError()
	filename := this.f.Name()
	idxfilename := this.idxf.Name()
	this.f.Close()
	this.idxf.Close()
	os.Remove(filename)
	os.Remove(idxfilename)
	return
}

func (this *Fdata) GetData(md5Bean *Md5Bean) (bs []byte, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	bs, err = ReadAt(this.f, int(md5Bean.Size), md5Bean.Offset)
	return
}

func (this *Fdata) Compact(chip int32) (finish bool) {
	defer catchError()
	if this.f.Name() == fm.currFileName || this.FileSize() > int64(chip*10) {
		return
	}
	//	fmt.Println("Compact:", this)
	bs, err := ioutil.ReadFile(this.idxf.Name())
	if err == nil {
		newfdata := fm._newFdata(false)
		length := len(bs) / 16
		finish = true
		for i := 0; i < length; i++ {
			md5key := bs[i*16 : (i+1)*16]
			mb, err := DBGetMd5Bean(md5key)
			if err == nil {
				bs, err := this.GetData(&mb)
				if err == nil {
					err = _AppendData(bs, newfdata)
				}
				if err != nil {
					finish = false
				}
			} else {
				fmt.Println("no md5key")
			}
			time.Sleep(10 * time.Millisecond)
		}
		if finish {
			fmt.Println("compact file: ", this.f.Name(), ">>>>>>", newfdata.f.Name())
			fmt.Println("compact size: ", this.FileSize(), ">>>>>>", newfdata.FileSize())
			this.CloseAndDelete()
		}
	}
	return
}

//--------------------------------------------------------------------------------------------------------------------

type Segment struct {
	Id       int64  //文件ID号
	Name     string //文件名
	FileType string //文件类型
	Md5      []byte //文件md5值
}

func NewSegment(name string, filetype string, md5 []byte) (s *Segment) {
	//	fmt.Println(name, " | ", filetype)
	s = new(Segment)
	s.Name = name
	s.FileType = filetype
	s.Md5 = md5
	return
}
func Bytes2Segment(bs []byte) (s Segment) {
	s = DecoderSegment(bs)
	return
}

//----------------------------------------------------------------------------------------------------------------------

type Md5Bean struct {
	Offset   int64  //文件所在位置
	Size     int32  //文件大小 字节
	FileName string //所在文件名
	QuoteNum int32  //引用数
	//	IsZip    bool   //是否zip压缩
	Sequence []byte //文件序号
}

func NewMd5Bean(offset int64, size int32, filename string, sequence []byte) (mb *Md5Bean) {
	mb = &Md5Bean{Offset: offset, Size: size, QuoteNum: 0, FileName: filename, Sequence: sequence}
	return
}

func Byte2Md5Bean(bs []byte) (md5 Md5Bean) {
	md5 = DecoderMd5(bs)
	return
}

func (this *Md5Bean) AddQuote() {
	atomic.AddInt32(&this.QuoteNum, 1)
}
func (this *Md5Bean) SubQuote() {
	atomic.AddInt32(&this.QuoteNum, -1)
}

//--------------------------------------------------------------------------------------------------------------------------

//type OctalBean struct {
//	Offset   int64  //文件所在位置
//	Size     int32  //文件大小
//	FileName string //所在文件名
//}
//func NewOctalBean(offset int64, size int32, filename string) (ob *OctalBean) {
//	mb = &OctalBean{Offset: offset, Size: size, FileName: filename}
//	return
//}
//func Byte2OctalBean(bs []byte) (ob OctalBean) {
//	ob = DecoderOctal(bs)
//	return
//}

//--------------------------------------------------------------------------------------------------------------------------
func DBPutSegment(key []byte, s Segment) (err error) {
	err = db.Put(key, EncodeSegment(s))
	return
}

func DBGetSegment(key []byte) (s Segment, err error) {
	var v []byte
	v, err = db.Get(key)
	if err == nil {
		s = DecoderSegment(v)
	}
	return
}

func DBPutMd5Bean(md5key []byte, md5Bean Md5Bean) (err error) {
	err = db.Put(md5key, EncodeMd5(md5Bean))
	return
}

func DBGetMd5Bean(md5key []byte) (md5 Md5Bean, err error) {
	var v []byte
	v, err = db.Get(md5key)
	if err == nil {
		md5 = DecoderMd5(v)
	}
	return
}

func DBDel(key []byte) (err error) {
	err = db.Del(key)
	return
}

func DBPut(key, value []byte) (err error) {
	err = db.Put(key, value)
	return
}

func DBExsit(key []byte) (b bool) {
	b = db.DBExist(key)
	return
}

//----------------------------------------------------------------------------------------------------------------------
func Append(f *os.File, bs []byte, offset int64) (n int, err error) {
	n, err = f.WriteAt(bs, offset)
	f.Sync()
	return
}
func Write(f *os.File, bs []byte) (n int, err error) {
	n, err = f.Write(bs)
	f.Sync()
	return
}

func ReadAt(f *os.File, byteInt int, offset int64) (bs []byte, err error) {
	bs = make([]byte, byteInt)
	_, err = f.ReadAt(bs, offset)
	return
}

//--------------------------------------------------------------------------------------------------------------------
func Hex2bytes(row int64) (bs []byte) {
	bs = make([]byte, 0)
	for i := 0; i < 8; i++ {
		r := row >> uint((7-i)*8)
		bs = append(bs, byte(r))
	}
	return
}

func Bytes2hex(bb []byte) (value int64) {
	value = int64(0x00000000)
	for i, b := range bb {
		ii := uint(b) << uint((7-i)*8)
		value = value | int64(ii)
	}
	return
}

func Octal2bytes(row int32) (bs []byte) {
	bs = make([]byte, 0)
	for i := 0; i < 4; i++ {
		r := row >> uint((3-i)*4)
		bs = append(bs, byte(r))
	}
	return
}

func Bytes2Octal(bb []byte) (value int32) {
	value = int32(0x0000)
	for i, b := range bb {
		ii := uint(b) << uint((3-i)*4)
		value = value | int32(ii)
	}
	return
}

//--------------------------------------------------------------------------------------------------------------------
func DecoderSegment(data []byte) (segment Segment) {
	var network bytes.Buffer
	_, er := network.Write(data)
	dec := gob.NewDecoder(&network)
	if er == nil {
		er = dec.Decode(&segment)
	}
	return
}

func EncodeSegment(segment Segment) (bs []byte) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(segment)
	bs = network.Bytes()
	return
}

//-------------------------------------------------
func DecoderMd5(data []byte) (md5 Md5Bean) {
	var network bytes.Buffer
	_, er := network.Write(data)
	dec := gob.NewDecoder(&network)
	if er == nil {
		er = dec.Decode(&md5)
	}
	return
}

func EncodeMd5(md5 Md5Bean) (bs []byte) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(md5)
	bs = network.Bytes()
	return
}

//-------------------------------------------------
//func DecoderOctal(data []byte) (ob OctalBean) {
//	var network bytes.Buffer
//	_, er := network.Write(data)
//	dec := gob.NewDecoder(&network)
//	if er == nil {
//		er = dec.Decode(&ob)
//		fmt.Println("ob===>", ob)
//	} else {
//		fmt.Println("er==>", er.Error())
//	}
//	return
//}

//func EncodeOctal(ob OctalBean) (bs []byte) {
//	var network bytes.Buffer
//	enc := gob.NewEncoder(&network)
//	enc.Encode(ob)
//	bs = network.Bytes()
//	return
//}

//--------------------------------------------------------------------------------------------------------------------

func MD5(data []byte) []byte {
	m := md5.New()
	m.Write(data)
	return m.Sum(nil)
}

//--------------------------------------------------------------------------------------------------------------------
func catchError(msg ...string) {
	if err := recover(); err != nil {
		if msg != nil {
			fmt.Println(strings.Join(msg, ","), err)
		}
	}
}

func _Ticker(second int, function func()) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	for {
		time.Sleep(time.Duration(second) * time.Second)
		function()
	}
}

func _compact() {
	catchError()
	m := db.GetIterLimit(_del_, fmt.Sprint(_del_, "z"))
	if m != nil {
		for k, v := range m {
			chip := Bytes2Octal(v)
			//			fmt.Println("scan compact key:", k, " ", chip)
			filename := strings.Replace(k, _del_, "", -1)
			fdata := fm.GetFdataByName(filename)
			if fdata.Compact(chip) {
				db.Del([]byte(k))
			}
		}
	}
}

func _fingerprint(bs []byte) (dest string) {
	ieee := crc32.NewIEEE()
	ieee.Write(bs)
	return hex.EncodeToString(ieee.Sum(nil))
}
