// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package stor

import (
	"bytes"
	"crypto/rand"
	"fmt"
	_ "net/http/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/wfs/sys"
)

func init() {
	sys.FileSize = 30 * sys.MB
	initStore()
	// go func ()  {
	// 	if err := http.ListenAndServe(":7000", nil); err != nil {
	// 		sys.FmtLog("debug  start failed:" + err.Error())
	// 	}
	// }()

}

var header = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}

func Test_stroge(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i <= 300000; i++ {
		wg.Add(1)
		func(i int) {
			defer wg.Done()
			bs := make([]byte, 1024)
			copy(bs[:10], header)
			rand.Read(bs[10:])
			if _, err := fe.append(fmt.Sprint("/aaa/b/c", i), bs, sys.CompressType); err != nil && !err.Equal(sys.ERR_EXSIT) {
				logging.Debug(i, ">>>", err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println(fe.delData("/aaa/b/c0"))

	var wg2 sync.WaitGroup
	for i := 0; i <= 300000; i++ {
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()
			if bs := fe.getData(fmt.Sprint("/aaa/b/c", i)); bs == nil {
				fmt.Println(i, ">>>>>", goutil.BytesToInt64(bs))
			}
		}(i)
	}
	wg2.Wait()

	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_stroge2(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i <= 1024; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bs := goutil.Int64ToBytes(int64(i))
			if _, err := fe.append(fmt.Sprint("/aaa/b/c", i), bs, sys.CompressType); err != nil {
				logging.Debug(i, ">>>", err)
			}
		}(i)
	}
	wg.Wait()

	// fmt.Println(fm.delData("/aaa/b/c0"))

	var wg2 sync.WaitGroup
	for i := 0; i <= 1024; i++ {
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()
			bs := fe.getData(fmt.Sprint("/aaa/b/c", i))
			if bs == nil {
				fmt.Println(i, ">>>>>", goutil.BytesToInt64(bs))
			}
		}(i)
	}
	wg2.Wait()
	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_stroge3(t *testing.T) {
	wps := fe.findLike("/aaa/b/c20000")
	for _, v := range wps {
		fmt.Println(v.Path, ",", time.Unix(0, v.Timestramp).Format("2006-01-02 15:04:05"), ",", len(v.Body))
	}
	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_del1(t *testing.T) {
	i := int64(0)
	for {
		if err := fe.delData(fmt.Sprint("/aaa/b/c", atomic.AddInt64(&i, 1))); err != nil {
			fmt.Println(err)
		}

		if i > 50000 {
			break
		}
	}
	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_stroge4(t *testing.T) {
	wps := fe.findLimit(5000, 10)
	for _, v := range wps {
		fmt.Println(v.Path, ",", time.Unix(0, v.Timestramp).Format("2006-01-02 15:04:05"), ",", len(v.Body))
	}
	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_get(t *testing.T) {
	var wg2 sync.WaitGroup
	go func() {
		defer wg2.Done()
		wg2.Add(1)
		fe.defragAndCover("4Tcrucsxofw")
	}()
	for i := 0; i <= 300000; i++ {
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()
			if bs := fe.getData(fmt.Sprint("/aaa/b/c", i)); bs == nil {
				fmt.Println(i, ">>>>>", goutil.BytesToInt64(bs))
			} else if !bytes.Equal(bs[:10], header) {
				fmt.Println("err>>>", bs[:10])
			}
		}(i)
	}
	wg2.Wait()
	fmt.Println(">>>>>>>>>>end>>>>>>>>>")
}

func Test_fragAnalysis(t *testing.T) {
	if fg, err := fe.fragAnalysis("CFvGjPLFBo4"); err == nil {
		fmt.Println(fg.Node, ",", fg.FileSize, ",", fg.ActualSize, ",", fg.RmSize)
	} else {
		fmt.Println(err)
	}
}

func Test_defrag(t *testing.T) {
	//BSWJ82iceuW , 1048576 , 1048520 , 100010
	fmt.Println(fe.defrag("AiVY4NJVFQ4"))
}

func Benchmark_Append(b *testing.B) {
	bsbase := make([]byte, 1<<10)
	rand.Read(bsbase)
	i := int64(0)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fe.append(fmt.Sprint("/aaa/b/c", atomic.AddInt64(&i, 1)), append(bsbase, goutil.Int16ToBytes(int16(goutil.RandId()))...), sys.CompressType)
		}
	})
}

func Benchmark_Get(b *testing.B) {
	i := int64(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fe.getData(fmt.Sprint("/aaa/b/c", atomic.AddInt64(&i, 1)))
		}
	})
}
