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

// StarImage 画像・動画を生成する上で変更したい値を定義しておく構造体
type StarImage struct {
	filename  string
	bgColor   color.Color
	starColor color.Color
}

func main() {
	files := []*StarImage{
		{
			filename:  "20211208/1.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{230, 0, 18, 255},
		},
		{
			filename:  "20211208/2.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{243, 152, 0, 255},
		},
		{
			filename:  "20211208/3.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{255, 251, 0, 255},
		},
		{
			filename:  "20211208/4.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{143, 195, 31, 255},
		},
		{
			filename:  "20211208/5.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{0, 153, 68, 255},
		},
		{
			filename:  "20211208/6.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{0, 158, 150, 255},
		},
		{
			filename:  "20211208/7.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{0, 160, 233, 255},
		},
		{
			filename:  "20211208/8.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{0, 104, 183, 255},
		},
		{
			filename:  "20211208/9.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{29, 32, 136, 255},
		},
		{
			filename:  "20211208/10.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{146, 7, 131, 255},
		},
		{
			filename:  "20211208/11.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{228, 0, 127, 255},
		},
		{
			filename:  "20211208/12.png",
			bgColor:   color.RGBA{0, 0, 64, 255},
			starColor: color.RGBA{229, 0, 79, 255},
		},
	}

	for _, f := range files {
		star := newStar(
			1000,
			1000,
			f.filename,
			f.bgColor,
			color.RGBA{0, 0, 0, 0},
			f.starColor,
		)
		star.drawStar()
	}

	// 各フレームの画像を GIF で読み込んで outGif を構築する
	outGif := &gif.GIF{}
	for _, f := range files {
		f, err := os.Open(f.filename)
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
	f, _ := os.OpenFile("20211208/out.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()

	err := gif.EncodeAll(f, outGif)
	if err != nil {
		log.Fatal(err)
	}
}
