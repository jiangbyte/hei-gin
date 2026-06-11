package utils

import (
	"log"
	"hei-gin/sdk/plugin"
	"hei-gin/sdk/config"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
)

var (
	_privateKey string
	_publicKey  string
	_sm2PrivKey *sm2.PrivateKey
	_sm2PubKey  *sm2.PublicKey
)

// Init initializes SM2 with the given private and public keys (hex strings).
// 0x prefix is automatically stripped and keys are lowercased.
func Init(privateKey, publicKey string) {
	_privateKey = strings.TrimSpace(strings.ToLower(privateKey))
	_privateKey = strings.TrimPrefix(_privateKey, "0x")
	_publicKey = strings.TrimSpace(strings.ToLower(publicKey))
	_publicKey = strings.TrimPrefix(_publicKey, "0x")

	if strings.TrimSpace(privateKey) == "" || strings.TrimSpace(publicKey) == "" {
		log.Println("[crypto] SM2: private or public key is empty, SM2 encryption will not be available")
		return
	}

	var err error
	_sm2PrivKey, err = hexToPrivateKey(_privateKey)
	if err != nil {
		log.Printf("[crypto] failed to parse SM2 private key: %v", err)
		return
	}
	_sm2PubKey, err = hexToPublicKey(_publicKey)
	if err != nil {
		log.Printf("[crypto] failed to parse SM2 public key: %v", err)
		return
	}
}

// hexToPrivateKey creates an sm2.PrivateKey from a raw hex string (64 hex chars = 32 bytes).
func hexToPrivateKey(hexStr string) (*sm2.PrivateKey, error) {
	dBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}
	d := new(big.Int).SetBytes(dBytes)
	if d.Sign() == 0 {
		return nil, fmt.Errorf("private key cannot be zero")
	}

	curve := sm2.P256Sm2()
	priv := new(sm2.PrivateKey)
	priv.D = d
	priv.PublicKey.Curve = curve
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	return priv, nil
}

// hexToPublicKey creates an sm2.PublicKey from a hex string.
// Supports both 64-byte (X||Y, no prefix) and 65-byte (04||X||Y) formats.
func hexToPublicKey(hexStr string) (*sm2.PublicKey, error) {
	pubBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, fmt.Errorf("invalid public key hex: %w", err)
	}

	curve := sm2.P256Sm2()
	offset := 0
	if len(pubBytes) > 0 && pubBytes[0] == 0x04 {
		offset = 1
	}

	keyLen := (len(pubBytes) - offset) / 2
	x := new(big.Int).SetBytes(pubBytes[offset : offset+keyLen])
	y := new(big.Int).SetBytes(pubBytes[offset+keyLen:])

	return &sm2.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// Encrypt encrypts plaintext using SM2 and returns the hex-encoded ciphertext in C1C2C3 format.
func Encrypt(plaintext string) string {
	if _sm2PubKey == nil {
		log.Println("[crypto] SM2 not initialized, Encrypt returning empty")
		return ""
	}
	ciphertext, err := sm2.Encrypt(_sm2PubKey, []byte(plaintext), rand.Reader, sm2.C1C2C3)
	if err != nil {
		log.Printf("[crypto] SM2 encryption failed: %v", err)
		return ""
	}
	return hex.EncodeToString(ciphertext)
}

// EncryptC1C3C2 encrypts plaintext using SM2 and returns the hex-encoded ciphertext
// reordered to C1C3C2 format (compatible with frontend sm-crypto cipherMode=1).
func EncryptC1C3C2(plaintext string) string {
	if _sm2PubKey == nil {
		log.Println("[crypto] SM2 not initialized, EncryptC1C3C2 returning empty")
		return ""
	}
	ciphertext, err := sm2.Encrypt(_sm2PubKey, []byte(plaintext), rand.Reader, sm2.C1C2C3)
	if err != nil {
		log.Printf("[crypto] SM2 encryption failed: %v", err)
		return ""
	}

	// ciphertext from sm2.Encrypt is C1||C2||C3
	// C1 = 65 bytes (04||x||y)
	// C3 = 32 bytes
	c1 := ciphertext[:65]
	c2 := ciphertext[65 : len(ciphertext)-32]
	c3 := ciphertext[len(ciphertext)-32:]

	// Reorder to C1||C3||C2
	result := make([]byte, 0, len(ciphertext))
	result = append(result, c1...)
	result = append(result, c3...)
	result = append(result, c2...)
	return hex.EncodeToString(result)
}

// Decrypt decrypts a hex-encoded SM2 ciphertext. Tries C1C2C3 format first,
// then falls back to C1C3C2 -> C1C2C3 conversion.
func Decrypt(ciphertext string) string {
	raw, err := decryptRaw(ciphertext)
	if err != nil {
		return ""
	}
	// Strip UTF-8 BOM if present
	if len(raw) >= 3 && raw[0] == 0xEF && raw[1] == 0xBB && raw[2] == 0xBF {
		raw = raw[3:]
	}
	return string(raw)
}

// decryptRaw performs the raw SM2 decryption, handling both C1C2C3 and C1C3C2 formats.
func decryptRaw(ciphertext string) ([]byte, error) {
	if _sm2PrivKey == nil {
		return nil, fmt.Errorf("SM2 has not been initialized. Please call Init() first.")
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext hex: %w", err)
	}

	// sm-crypto (cipherMode=1) outputs ciphertext without the 04 uncompressed
	// point prefix on C1. The tjfoc/gmsm library always strips the first byte
	// (expecting 04), so we must prepend it if missing.
	if len(data) > 0 && data[0] != 0x04 {
		data = append([]byte{0x04}, data...)
	}

	// Minimum valid length: 04(1) + X(32) + Y(32) + C3(32) = 97 bytes
	if len(data) < 97 {
		return nil, fmt.Errorf("ciphertext too short, invalid format")
	}

	// Attempt 1: Try C1C3C2 mode (frontend cipherMode=1)
	plaintext, err := sm2.Decrypt(_sm2PrivKey, data, sm2.C1C3C2)
	if err == nil && plaintext != nil {
		return plaintext, nil
	}

	// Attempt 2: Try C1C2C3 mode as fallback
	plaintext, err = sm2.Decrypt(_sm2PrivKey, data, sm2.C1C2C3)
	if err == nil && plaintext != nil {
		return plaintext, nil
	}

	return nil, fmt.Errorf("SM2 decryption failed, invalid ciphertext or wrong key")
}

// HashWithSalt computes the SM3 hash of data+salt and returns the hex-encoded result.
func HashWithSalt(data, salt string) string {
	h := sm3.New()
	h.Write([]byte(data + salt))
	return hex.EncodeToString(h.Sum(nil))
}

// GenSalt generates a random hex salt of the specified length (in hex characters).
func GenSalt(length int) string {
	if length <= 0 {
		return ""
	}
	// length is the desired hex string length, so we need length/2 bytes
	byteLen := length / 2
	if length%2 != 0 {
		byteLen++
	}
	bytes := make([]byte, byteLen)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// GenKeypair generates a new SM2 keypair.
// Returns (privateKeyHex, publicKeyHex).
func GenKeypair() (string, string) {
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		log.Printf("[crypto] failed to generate SM2 keypair: %v", err)
		return "", ""
	}

	// Private key: 32 bytes, zero-padded to 64 hex chars
	privHex := fmt.Sprintf("%064x", priv.D)

	// Public key: X||Y (64 bytes, no 04 prefix)
	pubBytes := make([]byte, 64)
	xBytes := bigIntTo32Bytes(priv.PublicKey.X)
	yBytes := bigIntTo32Bytes(priv.PublicKey.Y)
	copy(pubBytes[:32], xBytes)
	copy(pubBytes[32:], yBytes)
	pubHex := hex.EncodeToString(pubBytes)

	return privHex, pubHex
}

// GetPublicKey returns the initialized public key hex string.
func GetPublicKey() string {
	if _publicKey == "" {
		log.Println("[crypto] SM2 not initialized, GetPublicKey returning empty")
		return ""
	}
	return _publicKey
}

// bigIntTo32Bytes converts a big.Int to a fixed 32-byte slice.
func bigIntTo32Bytes(n *big.Int) []byte {
	b := n.Bytes()
	if len(b) >= 32 {
		return b[len(b)-32:]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}


// ---- plugin registration ----

type utilsPlugin struct{ plugin.NoopPlugin }

func (m *utilsPlugin) Name() string { return "utils" }

func (m *utilsPlugin) Init() error {
	Init(config.C.SM2.PrivateKey, config.C.SM2.PublicKey)
	return nil
}

func init() { plugin.Register(&utilsPlugin{}) }
