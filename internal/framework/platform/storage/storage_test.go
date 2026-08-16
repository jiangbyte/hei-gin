// internal/framework/platform/storage/storage_test.go
package storage

import "testing"

func TestEngineToProvider(t *testing.T) {
	cases := map[string]string{
		"MINIO":   ProviderMinIO,
		"RUSTFS":  ProviderRustFS,
		"ALIYUN":  ProviderOSS,
		"TENCENT": ProviderS3,
		"local":   "",
		"LOCAL":   "",
	}
	for in, want := range cases {
		if got := EngineToProvider(in); got != want {
			t.Fatalf("EngineToProvider(%q)=%q want %q", in, got, want)
		}
	}
}

func TestResolveProvider(t *testing.T) {
	if got := ResolveProvider("RUSTFS"); got != ProviderRustFS {
		t.Fatalf("got %q", got)
	}
	if got := ResolveProvider("minio"); got != ProviderMinIO {
		t.Fatalf("got %q", got)
	}
	if got := ResolveProvider("LOCAL"); got != "" {
		t.Fatalf("LOCAL should be rejected, got %q", got)
	}
	if !IsS3Compatible("oss") {
		t.Fatal("oss should be s3 compatible")
	}
	if IsS3Compatible("local") {
		t.Fatal("local must not be s3 compatible")
	}
}

func TestStripToObjectKey(t *testing.T) {
	cases := map[string]string{
		"uploads/a.txt":                "uploads/a.txt",
		"/api/v1/files/uploads/a.txt":  "uploads/a.txt",
		"api/v1/files/uploads/a.txt":   "uploads/a.txt",
		"https://cdn.example/x/y.png":  "x/y.png",
		"http://127.0.0.1:9000/vms/uploads/a.png": "uploads/a.png",
	}
	for in, want := range cases {
		if got := StripToObjectKey(in); got != want {
			t.Fatalf("StripToObjectKey(%q)=%q want %q", in, got, want)
		}
	}
}

func TestProviderConfigKeyPrefix(t *testing.T) {
	if got := ProviderConfigKeyPrefix("rustfs"); got != "STORAGE_RUSTFS" {
		t.Fatalf("got %q", got)
	}
	if got := ProviderConfigKeyPrefix("ALIYUN"); got != "STORAGE_ALIYUN" {
		t.Fatalf("got %q", got)
	}
	if got := ProviderConfigKeyPrefix("local"); got != "" {
		t.Fatalf("local prefix should be empty, got %q", got)
	}
}
