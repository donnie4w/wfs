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
