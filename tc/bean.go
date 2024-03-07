// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package tc

type AdminView struct {
	Show       string
	AdminUser  map[string]string
	Init       bool
	ShowCreate bool
}

type FilePage struct {
	CurrentNum  int
	ClientPort  int
	RevProxy    string
	CliProtocol string
	TotalNum    int
	FS          []*FileBean
}

type FileBean struct {
	Id   int
	Name string
	Time string
	Size int
}

type FragmentBean struct {
	Name         string
	Time         string
	Status       int
	FragmentSize int64
	FileSize     int64
}
