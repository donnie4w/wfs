package storge

import (
	"encoding/hex"
	"fmt"
	"hash/crc32"

	//	"io"
	"testing"
	"time"

	. "wfs/db"
)

func _Test_lockstring(t *testing.T) {
	bs := []byte{0, 0, 0, 0}
	bs1 := []byte{0, 0, 0, 2}
	go _lockstring(string(bs))
	go _lockstring(string(bs))
	go _lockstring(string(bs))
	go _lockstring(string(bs1))
	go _lockstring(string(bs1))
	time.Sleep(20 * time.Second)
}

func _lockstring(s string) {
	lockString.Lock(s)
	defer lockString.UnLock(s)
	fmt.Println(">>>>>", s)
	time.Sleep(3 * time.Second)
	fmt.Println("<<<<<", s)
}

func _Test_compact(t *testing.T) {
	db = NewDB("../data/fsdb", false)
	//	compact()
}

func Test_crc32(t *testing.T) {
	ieee := crc32.NewIEEE()
	ieee.Write([]byte("wuxiaodong0000000000000000"))
	fmt.Println(hex.EncodeToString(ieee.Sum(nil)))
}
