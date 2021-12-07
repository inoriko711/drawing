package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"

	"github.com/soniakeys/quant/median"
)

func main() {
	star := newStar(
		1000,
		1000,
		"20211208.png",
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 0, 0, 0},
		color.RGBA{255, 215, 0, 255},
	)
	star.drawStar()

	// 動画を作成する
	files := []string{
		"images/1.png",
		"images/2.png",
		"images/3.png",
		"images/4.png",
		"images/5.png",
		"images/6.png",
		"images/7.png",
		"images/8.png",
		"images/9.png",
		"images/10.png",
		"images/11.png",
		"images/10.png",
		"images/9.png",
		"images/8.png",
		"images/7.png",
		"images/6.png",
		"images/5.png",
		"images/4.png",
		"images/3.png",
		"images/2.png",
	}

	// 各フレームの画像を GIF で読み込んで outGif を構築する
	outGif := &gif.GIF{}
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer f.Close()

		// image.Imageへとデコード
		img, err := png.Decode(f)
		if err != nil {
			log.Fatal(err)
			return
		}
		q := median.Quantizer(256)
		p := q.Quantize(make(color.Palette, 0, 256), img)
		paletted := image.NewPaletted(img.Bounds(), p)
		draw.FloydSteinberg.Draw(paletted, img.Bounds(), img, image.Point{})

		outGif.Image = append(outGif.Image, paletted)
		outGif.Delay = append(outGif.Delay, 0)
	}

	// out.gif に保存する
	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}
