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
	"sort"
)

// パレットサイズと作成した画像を保存する場所
var (
	width    float64 = 1000
	height   float64 = 1000
	filename string  = "paintedStar2.png"
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
	X     int
	Y     int
	Color color.Color
}

func init() {
	centerX = width / 2.0
	centerY = height / 2.0
}

// 2点を受け取ってその間の座標を返す
func getCoordinates(x1, y1, x2, y2 float64) []Point {
	dx := x2 - x1
	dy := y2 - y1

	// 距離
	length := math.Sqrt(dx*dx + dy*dy)

	// ラジアン
	radian := math.Atan2(dy, dx)

	var points []Point
	for l := 0.0; l < length; l++ {
		// x座標とy座標を計算
		x := x1 + l*math.Cos(radian)
		y := y1 + l*math.Sin(radian)

		// ビットマップ外の点は描写しない
		if (x >= 0 && x < width) && (y >= 0 && y < height) {
			points = append(points, Point{x, y})
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
		fmt.Printf("外側：(x,y)=(%d,%d)\n", int(math.Round(x)), int(math.Round(y)))
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
		fmt.Printf("内側：(x,y)=(%d,%d)\n", int(math.Round(x)), int(math.Round(y)))
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

	sort.Slice(new, func(i, j int) bool { return new[i].X < new[j].X })
	return new
}

func convertPint2DrawPoint(points []Point, color color.Color) []DrawPoint {
	var drawPoints []DrawPoint
	for _, p := range points {
		drawPoints = append(drawPoints, DrawPoint{int(math.Round(p.X)), int(math.Round(p.Y)), color})
	}
	return deleteDuplicate(drawPoints)
}

func includePoint(x, y int, points []DrawPoint) bool {
	for _, p := range points {
		if x == p.X && y == p.Y {
			return true
		}
	}
	return false
}

func main() {
	m := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	pOutside := pentagonOutside()
	pInside := pentagonInside()

	// 星の外枠の座標を求める
	var starOutsidePoints []Point
	for i := 0; i < len(pOutside); i++ {
		starOutsidePoints = append(starOutsidePoints, getCoordinates(pOutside[i].X, pOutside[i].Y, pInside[i].X, pInside[i].Y)...)

		if i+1 == len(pOutside) {
			// pointsInsideの最後の点はpointsOutsideの最初の点につなげる
			starOutsidePoints = append(starOutsidePoints, getCoordinates(pInside[i].X, pInside[i].Y, pOutside[0].X, pOutside[0].Y)...)
		} else {
			starOutsidePoints = append(starOutsidePoints, getCoordinates(pInside[i].X, pInside[i].Y, pOutside[i+1].X, pOutside[i+1].Y)...)
		}
	}

	// 星の外枠の座標
	outsideDrawPoints := convertPint2DrawPoint(starOutsidePoints, linecolor)

	var drawPoint []DrawPoint
	for x := 0; x <= int(width); x++ {
		isStar := false
		beforePointColor := bgcolor

		for y := 0; y <= int(height); y++ {
			isLine := includePoint(x, y, outsideDrawPoints)
			// もしライン上の点だったらbeforePointColorをライン色に変えて次の点を確認する
			if isLine {
				isLine = true
				beforePointColor = linecolor
				drawPoint = append(drawPoint, DrawPoint{x, y, linecolor})
				continue
			}

			// 直前がラインのとき、isStarを反転させる
			if beforePointColor == linecolor {
				isStar = !isStar

				if x == int(width) || x == int(width)-1 { // TODO頂点処理
					isStar = false
				} else if x == 250 && y == 502 { // TODO 内側向き頂点処理
					isStar = true
				} else if x < int(centerX) && y <= int(centerY) {
					if includePoint(x+1, y+1, outsideDrawPoints) {
						isStar = false
					}
				} else if x < int(centerX) && y > int(centerY) {
					if includePoint(x+1, y-1, outsideDrawPoints) {
						isStar = false
					}
				}
			}

			if isStar {
				drawPoint = append(drawPoint, DrawPoint{x, y, starColor})
				beforePointColor = starColor
			} else {
				drawPoint = append(drawPoint, DrawPoint{x, y, bgcolor})
				beforePointColor = bgcolor
			}
		}
	}

	for _, p := range drawPoint {
		m.Set(p.X, p.Y, p.Color)
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
