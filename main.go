package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
)

// パレットサイズと作成した画像を保存する場所
var (
	width    int    = 500
	height   int    = 500
	filename string = "line.png"
)

// 背景色とラインカラー
var (
	bgcolor   color.Color = color.RGBA{255, 255, 255, 255}
	linecolor color.Color = color.RGBA{0, 0, 0, 0}
)

// 2点を受け取って線を引く
func drawLine(m *image.RGBA, x1, y1, x2, y2 float64) {
	dx := x2 - x1
	dy := y2 - y1

	// 距離
	length := math.Sqrt(dx*dx + dy*dy)
	fmt.Println("length: ", length)

	// ラジアン
	radian := math.Atan2(dy, dx)
	fmt.Println("radian: ", radian)
	var deg = radian * 180 / math.Pi
	fmt.Printf("角度：%f\n", deg)

	for l := 0.0; l < length; l++ {
		// x座標とy座標を計算
		x := x1 + l*math.Cos(radian)
		y := y1 + l*math.Sin(radian)
		fmt.Printf("(x,y): (%f,%f)\n", x, y)

		// ビットマップ外の点は描写しない
		if (x >= 0 && int(x) < width) && (y >= 0 && int(y) < height) {
			m.Set(int(x), int(y), linecolor)
		}
	}
}

func main() {
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	drawLine(m, 10, 10, 10, 400)
	drawLine(m, 10, 400, 400, 400)
	drawLine(m, 400, 400, 400, 10)
	drawLine(m, 400, 10, 10, 10)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, m)
	if err != nil {
		log.Fatal(err)
	}
}
