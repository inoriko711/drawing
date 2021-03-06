package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	m := image.NewRGBA(image.Rect(0, 0, 11, 13))
	draw.Draw(m, m.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// 始点(1,6)から終点(9,6)までの間の点全てに色を塗る
	for x := 1; x < 10; x++ {
		m.Set(x, 6, color.RGBA{0, 0, 0, 0})
	}

	// ファイルに出力する
	file, err := os.Create("line.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, m)
	if err != nil {
		log.Fatal(err)
	}
}
