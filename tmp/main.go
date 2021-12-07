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
	m.Set(5, 6, color.RGBA{0, 0, 0, 0})
	file, err := os.Create("point.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, m)
	if err != nil {
		log.Fatal(err)
	}
}
