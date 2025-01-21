package gcaptha

import (
	"crypto/rand"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math/big"
	"os"
)

type solidImg struct {
	w, h, offset int
}

// 生成信息
func (this solidImg) GenCode() map[string]any {
	return map[string]any{
		"width":  this.w,
		"height": this.h,
		"tk":     "asdf",
	}
}

// 验证
func (this solidImg) Check(tk, code string) bool {
	return false
}

// 生成图片
func (this solidImg) Encode(out io.Writer) {
	bh := this.w / 8
	bw := bh + bh/4
	img := image.NewRGBA(image.Rect(0, 0, this.w, this.h*2))
	bg := this.bg()
	num, err := rand.Int(rand.Reader, big.NewInt(int64(this.h-10-bw)))
	if err != nil {
		fmt.Printf("err %s\n", err.Error())
	}
	y := num.Int64() + 5
	// fmt.Println(y, this.h+int(y))
	draw.Draw(img, image.Rect(0, 0, this.w, this.h), bg, image.Point{}, draw.Src)
	block := this.block(true)
	x := bw + 5

	// 创建一个遮罩图像
	mask := image.NewNRGBA(image.Rect(0, 0, bw, bh))
	draw.Draw(mask, mask.Bounds(), image.NewUniform(color.RGBA{R: 0, G: 0, A: 220}), image.Point{}, draw.Src)
	draw.DrawMask(img, image.Rect(5+this.offset, int(y), x+this.offset, int(y)+bh), mask, image.Point{}, block, image.Point{}, draw.Over)

	subimg := bg.SubImage(image.Rect(5+this.offset, int(y), x+this.offset, int(y)+bw))
	mask2 := image.NewAlpha(image.Rect(0, 0, bw, bh))
	draw.Draw(mask2, mask2.Bounds(), image.NewUniform(color.RGBA{R: 0, G: 0, A: 220}), image.Point{}, draw.Src)
	// draw.Draw(img, image.Rect(5, this.h+int(y), x, this.h+int(y)+bw), block, image.Point{}, draw.Over)
	if simg, ok := subimg.(*image.RGBA); ok {
		// draw.Draw(img, image.Rect(5, this.h+int(y), x, this.h+int(y)+bw), bg, image.Point{X: 5 + this.offset, Y: int(y)}, draw.Over)
		draw.DrawMask(img, image.Rect(5, this.h+int(y), x, this.h+int(y)+bw), simg, image.Point{X: 5 + this.offset, Y: int(y)}, block, image.Point{}, draw.Over)
		// fmt.Println("生成字图片")

	}
	block2 := this.block(false)
	draw.Draw(img, image.Rect(5, this.h+int(y), x, this.h+int(y)+bw), block2, image.Point{}, draw.Over)

	png.Encode(out, img)
}
func (this solidImg) bg() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, this.w, this.h))
	f, _ := os.Open("1.png")
	defer f.Close()
	simg, _ := png.Decode(f)
	draw.Draw(img, img.Bounds(), simg, image.Point{}, draw.Over)
	return img
}
func (this solidImg) block(ismark bool) image.Image {
	w := this.w / 8
	r := w / 4
	img := image.NewRGBA(image.Rect(0, 0, w+r, w))

	if ismark {
		draw.Draw(img, img.Bounds(), image.NewUniform(color.Alpha{255}), image.Point{}, draw.Over)
	} else {
		for i := 0; i <= w; i++ {
			img.Set(i, 0, color.Black)
			img.Set(i, w-1, color.Black)

			img.Set(i, 1, color.Black)
			img.Set(i, w-2, color.Black)

			img.Set(i, 2, color.Black)
			img.Set(i, w-3, color.Black)
		}
		for i := 0; i <= r; i++ {
			img.Set(0, i, color.Black)
			img.Set(w-1, i, color.Black)
			img.Set(0, w-i, color.Black)
			img.Set(w-1, w-i, color.Black)

			img.Set(1, i, color.Black)
			img.Set(w-2, i, color.Black)
			img.Set(1, w-i, color.Black)
			img.Set(w-2, w-i, color.Black)

			img.Set(2, i, color.Black)
			img.Set(w-3, i, color.Black)
			img.Set(2, w-i, color.Black)
			img.Set(w-3, w-i, color.Black)

		}

	}

	y := w / 2
	x1 := 0
	x2 := w
	r2 := r * r
	for xi := 0; xi <= r; xi++ {
		for yi := 0; yi <= 2*r; yi++ {
			yy := y - r + yi
			yr := yy - y
			xr := xi - x1
			rr := (xr * xr) + (yr * yr)
			if !ismark && rr <= r2 {
				img.Set(xi, yy, color.Black)
			} else if rr <= r2 {
				img.Set(xi, yy, color.Alpha{0})
			}
		}
	}
	for xi := 0; xi <= r; xi++ {
		for yi := 0; yi <= w; yi++ {
			yr := yi - y
			xx := x2 + xi
			xr := xx - x2
			rr := (xr * xr) + (yr * yr)
			if !ismark && rr <= r2 {
				img.Set(xx, yi, color.Black)
			} else if rr >= r2 {
				img.Set(xx, yi, color.Alpha{0})
			}
		}
	}
	if !ismark {
		r = r - 3
		r2 := r * r
		for xi := 0; xi <= r; xi++ {
			for yi := 0; yi <= 2*r; yi++ {
				yy := y - r + yi
				yr := yy - y
				xr := xi - x1
				rr := (xr * xr) + (yr * yr)
				if rr <= r2 {
					img.Set(xi, yy, color.Alpha{0})
				}
			}
		}
		for xi := 0; xi <= r; xi++ {
			for yi := 0; yi <= w; yi++ {
				yr := yi - y
				xx := x2 + xi
				xr := xx - x2
				rr := (xr * xr) + (yr * yr)
				if rr <= r2 {
					img.Set(xx, yi, color.Alpha{0})
				}
			}
		}
	}
	return img
}
