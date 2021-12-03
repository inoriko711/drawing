package main

import (
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
	width    float64 = 1000
	height   float64 = 1000
	filename string  = "paintedStar.png"
	// 中心のx,y座標
	centerX float64
	centerY float64
)

// 背景色とラインカラー
var (
	bgcolor   color.Color = color.RGBA{255, 255, 255, 255}
	linecolor color.Color = color.RGBA{0, 0, 0, 0}
	starColor color.Color = color.RGBA{255, 215, 0, 255}
)

type Point struct {
	X float64
	Y float64
}

type DrawPoint struct {
	X int
	Y int
}

func init() {
	centerX = width / 2.0
	centerY = height / 2.0
}

// 2点を受け取ってその間の座標を返す
func getCoordinates(x1, y1, x2, y2 float64) []DrawPoint {
	dx := x2 - x1
	dy := y2 - y1

	// 距離
	length := math.Sqrt(dx*dx + dy*dy)

	// ラジアン
	radian := math.Atan2(dy, dx)

	var points []DrawPoint
	for l := 0.0; l < length; l++ {
		// x座標とy座標を計算
		x := x1 + l*math.Cos(radian)
		y := y1 + l*math.Sin(radian)

		// ビットマップ外の点は描写しない
		if (x >= 0 && x < width) && (y >= 0 && y < height) {
			points = append(points, DrawPoint{X: int(x), Y: int(y)})
		}
	}

	return points
}

// 星の外側の頂点座標を求める
func pentagonOutside() []Point {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := width / 2.0

	// 中心から頂点への直線同士が成す角度
	radian := math.Pi * 2 / vertexNum

	var points []Point
	for i := 0.0; i < vertexNum; i++ {
		x := centerX + R*math.Cos(radian*i)
		y := centerY + R*math.Sin(radian*i)
		points = append(points, Point{x, y})
	}

	return points
}

// 星の内側の頂点座標を求める
func pentagonInside() []Point {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := width / 4.0

	// 外側の1/2の角度を求める
	radian := math.Pi / vertexNum

	var points []Point
	for i := 1.0; i < vertexNum*2; i = i + 2 { // 外側と角度をずらしつつ幅を揃える
		x := centerX + R*math.Cos(radian*i)
		y := centerY + R*math.Sin(radian*i)
		points = append(points, Point{x, y})
	}

	return points
}

// Sliceの重複を削除する
func deleteDuplicate(old []DrawPoint) []DrawPoint {
	m := make(map[DrawPoint]bool)
	var new []DrawPoint

	for _, point := range old {
		if !m[point] {
			m[point] = true
			new = append(new, point)
		}
	}
	return new
}

func main() {
	m := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	pOutside := pentagonOutside()
	pInside := pentagonInside()

	// 星の外枠の座標を求める
	var starOutsidePoints []DrawPoint
	for i := 0; i < len(pOutside); i++ {
		starOutsidePoints = append(starOutsidePoints, getCoordinates(pOutside[i].X, pOutside[i].Y, pInside[i].X, pInside[i].Y)...)

		if i+1 == len(pOutside) {
			// pointsInsideの最後の点はpointsOutsideの最初の点につなげる
			starOutsidePoints = append(starOutsidePoints, getCoordinates(pInside[i].X, pInside[i].Y, pOutside[0].X, pOutside[0].Y)...)
		} else {
			starOutsidePoints = append(starOutsidePoints, getCoordinates(pInside[i].X, pInside[i].Y, pOutside[i+1].X, pOutside[i+1].Y)...)
		}
	}
	starOutsidePoints = deleteDuplicate(starOutsidePoints)

	// 星の内側の座標を求める
	var starInsidePoints []DrawPoint
	for _, p := range starOutsidePoints {
		starInsidePoints = append(starInsidePoints, getCoordinates(centerX, centerY, float64(p.X), float64(p.Y))...)
	}
	starInsidePoints = deleteDuplicate(starInsidePoints)

	for _, points := range starInsidePoints {
		m.Set(points.X, points.Y, starColor)
	}
	for _, points := range starOutsidePoints {
		m.Set(points.X, points.Y, linecolor)
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
