package gcaptha

import (
	"crypto/rand"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/big"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// 生成随机字符串
func GenCodeStr(slen int) string {
	const letters = "abcdefghjklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, slen)
	maxint := len(letters)
	for i := 0; i < slen; i++ {
		result[i] = letters[randNum(maxint)]
	}
	return string(result)
}

// 生成数学运算验证码
func GenMathStr() (drawStr, code string) {
	const ops = "+×"
	op := ops[randNum(2)]
	if op == '+' {
		s1 := randNum(50)
		s2 := randNum(50)
		return fmt.Sprintf("%d + %d =", s1, s2), fmt.Sprintf("%d", s1+s2)
	} else {
		s1 := randNum(10)
		s2 := randNum(10)
		return fmt.Sprintf("%d × %d", s1, s2), fmt.Sprintf("%d", s1*s2)
	}

}

// 生成随机数
func randNum(max int) int64 {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return num.Int64()
}

// 渲染字符验证码图片
func DrawCapthaStrImg(out io.Writer, w, h int, code string) {
	f, e := opentype.Parse(fontData)
	if e != nil {
		panic(e)
	}
	size := 300
	slen := len([]rune(code))
	aw := w / slen
	if size > h {
		size = h - 2
	}
	if size > aw {
		size = aw - 1
	}
	face, e := opentype.NewFace(f, &opentype.FaceOptions{
		Hinting: font.HintingNone,
		DPI:     72,
		Size:    float64(size),
	})
	if e != nil {
		panic(e)
	}
	box, _ := font.BoundString(face, code)
	boxW := (box.Max.X - box.Min.X).Ceil()
	boxH := (box.Max.Y - box.Min.Y).Ceil()
	img := image.NewGray(image.Rect(0, 0, w, h))
	d := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P((w-boxW)/2, (h+boxH)/2),
	}
	// fmt.Printf("[%d,%d]字体size %d 字符串宽度 %d dot(%d,%d) box:%v\n", w, h, size, size*slen, (w-boxW)/2, (h+boxH)/2, box)
	d.DrawString(code)

	png.Encode(out, img)

}
