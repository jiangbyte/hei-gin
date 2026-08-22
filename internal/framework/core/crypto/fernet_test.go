package crypto

import "testing"

func TestFernetRoundTrip(t *testing.T) {
	// url-safe base64 32-byte key (dev sample from config.example.yaml)
	key := "XV1rJ-UPAbWeYjprihKNS3ZCCHdBuVbIc0WXmYc70ck="
	codec, err := NewFernet(key)
	if err != nil {
		t.Fatal(err)
	}
	plain := "hei-gin-secret"
	token, err := codec.Encrypt(plain)
	if err != nil {
		t.Fatal(err)
	}
	got, err := codec.Decrypt(token)
	if err != nil {
		t.Fatal(err)
	}
	if got != plain {
		t.Fatalf("want %q got %q", plain, got)
	}
}
