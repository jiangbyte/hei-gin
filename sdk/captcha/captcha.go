package captcha

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hei-gin/sdk/infra/db"
	"hei-gin/sdk/kernel/plugin"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Character set for captcha codes (no I, O, 0, 1 to avoid ambiguity).
const captchaChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// CaptchaResult represents the captcha generation result.
type CaptchaResult struct {
	CaptchaBase64 string `json:"captcha_base64"`
	CaptchaID     string `json:"captcha_id"`
	CaptchaCode   string `json:"captcha_code"`
}

// CaptchaService provides captcha generation and verification backed by Redis.
type CaptchaService struct {
	prefix string
	redis  *redis.Client
}

// NewCaptchaService creates a new CaptchaService with the given Redis key prefix.
func NewCaptchaService(prefix string) *CaptchaService {
	return &CaptchaService{prefix: prefix}
}

// Init sets the Redis client for the service.
func (s *CaptchaService) Init(redis *redis.Client) {
	s.redis = redis
}

// GetCaptcha generates a captcha image and returns the base64-encoded image, ID, and code.
// The code is stored in Redis with a 300-second TTL.
func (s *CaptchaService) GetCaptcha() (*CaptchaResult, error) {
	captchaID := uuid.New().String()

	// Generate 4-character code from the safe character set
	code := make([]byte, 4)
	for i := 0; i < 4; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(captchaChars))))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random character: %w", err)
		}
		code[i] = captchaChars[n.Int64()]
	}
	codeStr := string(code)

	// Create white background image (100x38)
	img := image.NewRGBA(image.Rect(0, 0, 100, 38))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Get font face (TrueType preferred, basicfont fallback)
	face := getFontFace(24)

	// Draw each character with random position and color
	for i, ch := range codeStr {
		x := fixed.I(10 + i*22)
		y := fixed.I(randomInt(22, 32))

		d := &font.Drawer{
			Dst: img,
			Src: image.NewUniform(color.RGBA{
				R: uint8(randomInt(0, 100)),
				G: uint8(randomInt(0, 100)),
				B: uint8(randomInt(0, 100)),
				A: 255,
			}),
			Face: face,
			Dot:  fixed.Point26_6{X: x, Y: y},
		}
		d.DrawString(string(ch))
	}

	// Draw 10 random noise lines
	for range 10 {
		x1 := randomInt(0, 100)
		y1 := randomInt(0, 38)
		x2 := randomInt(0, 100)
		y2 := randomInt(0, 38)
		drawLine(img, x1, y1, x2, y2, color.RGBA{
			R: uint8(randomInt(0, 200)),
			G: uint8(randomInt(0, 200)),
			B: uint8(randomInt(0, 200)),
			A: 255,
		})
	}

	// Encode to PNG and base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode captcha image: %w", err)
	}
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Store code in Redis with 300s TTL
	if s.redis != nil {
		err := s.redis.Set(context.TODO(), s.prefix+captchaID, codeStr, 300*time.Second).Err()
		if err != nil {
			log.Printf("[CAPTCHA] failed to store captcha code in Redis: %v", err)
		}
	}

	return &CaptchaResult{
		CaptchaBase64: "data:image/png;base64," + imageBase64,
		CaptchaID:     captchaID,
	}, nil
}

// CheckCaptcha verifies the captcha code against the value stored in Redis.
// The stored code is deleted after successful verification (one-time use).
func (s *CaptchaService) CheckCaptcha(id, code string) error {
	if id == "" || code == "" {
		return fmt.Errorf("验证码ID或验证码内容不能为空")
	}

	if s.redis == nil {
		log.Printf("[CAPTCHA] WARNING: Redis is not initialized, captcha check bypassed for id=%s", id)
		return nil
	}

	ctx := context.TODO()
	key := s.prefix + id
	storedCode, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("验证码已过期或无效")
	}
	if err != nil {
		return fmt.Errorf("failed to check captcha: %w", err)
	}

	// Case-insensitive comparison, trimming whitespace
	if !strings.EqualFold(strings.TrimSpace(code), storedCode) {
		return fmt.Errorf("验证码错误")
	}

	// Delete after successful verification
	if err := s.redis.Del(ctx, key).Err(); err != nil {
		log.Printf("[CAPTCHA] failed to delete captcha key %s: %v", key, err)
	}

	return nil
}

// getFontFace attempts to load a TrueType font from common system paths.
// Falls back to basicfont.Face7x13 if no TrueType font is available.
func getFontFace(size float64) font.Face {
	// Common TrueType font paths across platforms (TTF only, no TTC).
	fontPaths := []string{
		"C:\\Windows\\Fonts\\arial.ttf",
		"C:\\Windows\\Fonts\\Arial.ttf",
		"/usr/share/fonts/truetype/msttcorefonts/Arial.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/Library/Fonts/Arial.ttf",
	}

	var fontData []byte
	var readErr error
	for _, path := range fontPaths {
		fontData, readErr = os.ReadFile(path)
		if readErr == nil {
			break
		}
	}

	if readErr == nil && len(fontData) > 0 {
		f, parseErr := sfnt.Parse(fontData)
		if parseErr == nil {
			face, faceErr := opentype.NewFace(f, &opentype.FaceOptions{
				Size:    size,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if faceErr == nil {
				return face
			}
		}
	}

	// Fallback: Go's built-in basic bitmap font (7x13 fixed-width).
	return basicfont.Face7x13
}

// drawLine draws a line on the image using Bresenham's line algorithm.
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	var sx, sy int
	if x2 > x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y2 > y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy
	for {
		if x1 >= 0 && x1 < img.Bounds().Max.X && y1 >= 0 && y1 < img.Bounds().Max.Y {
			img.Set(x1, y1, c)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// randomInt returns a random integer in [min, max].
func randomInt(min, max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return min
	}
	return min + int(n.Int64())
}

// Package-level exported captcha service instances for business and consumer.
var (
	BCaptcha = NewCaptchaService("BUSINESS:captcha:")
	CCaptcha = NewCaptchaService("CONSUMER:captcha:")
)

// ---- plugin registration ----

type captchaPlugin struct{ plugin.NoopPlugin }

var registerOnce sync.Once

func (m *captchaPlugin) Name() string { return "captcha" }

func (m *captchaPlugin) Init() error {
	BCaptcha.Init(db.Redis)
	CCaptcha.Init(db.Redis)
	return nil
}

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&captchaPlugin{})
	})
}
