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
	width    float64 = 500
	height   float64 = 500
	filename string  = "line.png"
)

// 背景色とラインカラー
var (
	bgcolor   color.Color = color.RGBA{255, 255, 255, 255}
	linecolor color.Color = color.RGBA{0, 0, 0, 0}
)

type Points struct {
	X float64
	Y float64
}

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
		// fmt.Printf("(x,y): (%f,%f)\n", x, y)

		// ビットマップ外の点は描写しない
		if (x >= 0 && x < width) && (y >= 0 && y < height) {
			m.Set(int(x), int(y), linecolor)
		}
	}
}

func pentagonOutside() []Points {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := width / 2.0
	radian := math.Pi * 2 / vertexNum

	// 中心のx,y座標
	centerX := width / 2.0
	centerY := height / 2.0

	var points []Points
	for i := 0.0; i < vertexNum; i++ {
		x := centerX + R*math.Cos(radian*i)
		y := centerY + R*math.Sin(radian*i)
		points = append(points, Points{x, y})
	}

	return points
}

func pentagonInside() []Points {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := width / 4.0
	radian := math.Pi / vertexNum

	// 中心のx,y座標
	centerX := width / 2.0
	centerY := height / 2.0

	var points []Points
	for i := 1.0; i < vertexNum*2; i = i + 2 {
		x := centerX + R*math.Cos(radian*i)
		y := centerY + R*math.Sin(radian*i)
		points = append(points, Points{x, y})
	}

	return points
}
func main() {
	m := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	pointsOutside := pentagonOutside()
	for _, p := range pointsOutside {
		drawLine(m, 250, 250, p.X, p.Y)
	}

	pointsInside := pentagonInside()
	for _, p := range pointsInside {
		drawLine(m, 250, 250, p.X, p.Y)
	}
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
