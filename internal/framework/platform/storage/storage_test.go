package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestNewLocal_NoEagerDir 构造本地存储不应立即创建根目录（避免启动产生空 storage/）。
func TestNewLocal_NoEagerDir(t *testing.T) {
	root := filepath.Join(t.TempDir(), "storage")
	local := NewLocal(root, "/api/v1/files", "")
	if _, err := os.Stat(local.Root()); !os.IsNotExist(err) {
		t.Fatalf("root dir should NOT exist after NewLocal, stat err = %v", err)
	}
	// Put 写入后应惰性创建目录
	if _, err := local.Put(context.Background(), "uploads/2024/01/01/x.txt", bytes.NewBufferString("hi"), 2, "text/plain"); err != nil {
		t.Fatalf("put: %v", err)
	}
	if st, err := os.Stat(filepath.Join(root, "uploads")); err != nil || !st.IsDir() {
		t.Fatalf("dir should exist after Put, stat=%v err=%v", st, err)
	}
}
