// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package stor

import (
	goutil "github.com/donnie4w/gofer/util"
)

var (
	VERSION_       = []byte{127, 1}
	ENDOFFSET_     = []byte{1}
	APPENDLOCK_    = []byte{2}
	RESETMMAPLOCK_ = []byte{3}
	OPENMMAPLOCK_  = []byte{4}
	PATH_PRE       = []byte{0, 0}
	CURRENT        = append([]byte{6}, goutil.Int64ToBytes(1<<50)...)
	SEQ            = append([]byte{7}, goutil.Int64ToBytes(1<<51)...)
	PATH_SEQ       = append([]byte{8}, goutil.Int64ToBytes(1<<52)...)
	COUNT          = append([]byte{9}, goutil.Int64ToBytes(1<<53)...)
)
