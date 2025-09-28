// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package main

import (
	_ "github.com/donnie4w/wfs/keystore"
	_ "github.com/donnie4w/wfs/level1"
	_ "github.com/donnie4w/wfs/stor"
	. "github.com/donnie4w/wfs/sys"
	_ "github.com/donnie4w/wfs/tc"
)

func main() {
	Wfs.Serve()
}
