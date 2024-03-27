// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stor

import (
	"runtime"
	"time"

	"github.com/donnie4w/wfs/sys"
	"github.com/shirou/gopsutil/v3/mem"
)

type Ram struct {
	OsUsedMB    uint64
	OsTotalMB   uint64
	UserdMB     uint64
	UsedPercent float64
}

var ram = newRam()

func newRam() (r *Ram) {
	r = &Ram{}
	if u, err := mem.VirtualMemory(); err == nil {
		r.OsTotalMB = u.Total
		r.OsUsedMB = u.Used
		r.UsedPercent = u.UsedPercent
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	r.UserdMB = m.TotalAlloc
	return
}

func restrict() bool {
	return int(ram.UsedPercent) > sys.Restrict
}

func tasklimit() {
	for over := 100; over > 0 && restrict(); over-- {
		<-time.After(time.Second / time.Duration(over))
	}
}

func storTk() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			ram = newRam()
		}
	}
}
