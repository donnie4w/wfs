package main

import "testing"

func Test_main(t *testing.T) {
	leveldbDir := "../wfsdata/wfsdb"
	dbFile := "../wfsdata/wfsdb/wfs.db"
	t.Logf("开始迁移: (%s) -> (%s)", leveldbDir, dbFile)
	if err := migrateData(leveldbDir, dbFile); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
}
