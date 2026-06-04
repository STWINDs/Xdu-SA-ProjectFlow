package captcha

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// GenerateCaptcha creates a simple math captcha image (e.g. "3 + 5 = ?").
// It returns a unique id, a base64-encoded PNG image, the correct answer string, and any error.
func GenerateCaptcha() (id string, b64Image string, answer string, err error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	a := rng.Intn(20) + 1 // 1-20
	b := rng.Intn(20) + 1 // 1-20
	op := rng.Intn(2)     // 0: addition, 1: subtraction

	var opStr string
	var result int
	if op == 0 {
		opStr = "+"
		result = a + b
	} else {
		// Ensure non-negative result by swapping if needed.
		if a < b {
			a, b = b, a
		}
		opStr = "-"
		result = a - b
	}

	question := fmt.Sprintf("%d %s %d = ?", a, opStr, b)

	width := 200
	height := 60
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with a light color.
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Add random noise dots to make OCR harder.
	for i := 0; i < 100; i++ {
		x := rng.Intn(width)
		y := rng.Intn(height)
		noiseColor := color.RGBA{
			uint8(rng.Intn(128)),
			uint8(rng.Intn(128)),
			uint8(rng.Intn(128)),
			255,
		}
		img.Set(x, y, noiseColor)
	}

	// Draw random interfering lines.
	for i := 0; i < 3; i++ {
		x1 := rng.Intn(width)
		y1 := rng.Intn(height)
		x2 := rng.Intn(width)
		y2 := rng.Intn(height)
		lineColor := color.RGBA{
			uint8(rng.Intn(100) + 100),
			uint8(rng.Intn(100) + 100),
			uint8(rng.Intn(100) + 100),
			255,
		}
		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// Draw the math question text.
	textColor := color.RGBA{20, 50, 120, 255}
	point := fixed.Point26_6{
		X: fixed.Int26_6(15 * 64),
		Y: fixed.Int26_6(38 * 64),
	}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(question)

	// Encode image to PNG and then to base64.
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return "", "", "", err
	}

	id = fmt.Sprintf("captcha_%d", time.Now().UnixNano())
	b64Image = base64.StdEncoding.EncodeToString(buf.Bytes())
	answer = strconv.Itoa(result)

	return id, b64Image, answer, nil
}

// VerifyCaptcha checks the provided answer against the value stored in Redis under the captcha id.
// The stored value should have been set with a 5-minute TTL by the caller (use StoreCaptcha).
func VerifyCaptcha(rdb *redis.Client, id, answer string) bool {
	if rdb == nil {
		return false
	}
	key := "captcha:" + id
	stored, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if stored == answer {
		// Consume the captcha after successful verification.
		_ = rdb.Del(context.Background(), key)
		return true
	}
	return false
}

// StoreCaptcha stores the captcha answer in Redis with a 5-minute TTL.
// Call this after GenerateCaptcha to persist the captcha for later verification.
func StoreCaptcha(rdb *redis.Client, id, answer string) error {
	if rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	key := "captcha:" + id
	return rdb.Set(context.Background(), key, answer, 5*time.Minute).Err()
}

// drawLine is a simple Bresenham line drawing helper for captcha noise.
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy
	for {
		if x1 >= 0 && x1 < img.Bounds().Dx() && y1 >= 0 && y1 < img.Bounds().Dy() {
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
