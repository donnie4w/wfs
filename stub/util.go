// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stub

import "google.golang.org/protobuf/proto"

func (t *SnapshotBean) Copy(key, value []byte) {
	if key != nil {
		t.Key = make([]byte, len(key))
		copy(t.Key, key)
	}
	if value != nil {
		t.Value = make([]byte, len(value))
		copy(t.Value, value)
	}
}

func (t *SnapshotBean) ToBytes() ([]byte, error) {
	return proto.Marshal(t)
}

func BytesToSnapshotBean(bs []byte) (_r *SnapshotBean, err error) {
	_r = new(SnapshotBean)
	err = proto.Unmarshal(bs, _r)
	return
}

func (t *SnapshotBeans) ToBytes() ([]byte, error) {
	return proto.Marshal(t)
}

func BytesToSnapshotBeans(bs []byte) (_r *SnapshotBeans, err error) {
	_r = new(SnapshotBeans)
	err = proto.Unmarshal(bs, _r)
	return
}

func (t *SnapshotFile) ToBytes() ([]byte, error) {
	return proto.Marshal(t)
}

func BytesToSnapshotFile(bs []byte) (_r *SnapshotFile, err error) {
	_r = new(SnapshotFile)
	err = proto.Unmarshal(bs, _r)
	return
}