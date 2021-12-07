package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"sort"

	"github.com/soniakeys/quant/median"
)

// パレットサイズと作成した画像を保存する場所
var (
	width    float64 = 1000
	height   float64 = 1000
	filename string  = "paintedStarAtNight.png"
	// 中心のx,y座標
	centerX float64
	centerY float64
)

// 背景色とラインカラー
var (
	bgcolor   color.Color = color.RGBA{0, 0, 64, 255}
	linecolor color.Color = color.RGBA{0, 0, 0, 0}
	// -25,-21,+6
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

type LinePoint struct {
	FromY int
	ToY   int
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

	sort.Slice(new, func(i, j int) bool {
		if new[i].X == new[j].X {
			return new[i].Y < new[j].Y
		}
		return new[i].X < new[j].X
	})

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

		// 直前のマスがライン上かどうか保持する変数
		isLine := false
		yTo := 0
		yFrom := 0
		var linePoints []LinePoint

		// 星の枠線のy座標を求める
		for y := 0; y <= int(height); y++ {
			if !isLine && includePoint(x, y, outsideDrawPoints) {
				yFrom = y
				yTo = y
				isLine = true
			} else if isLine && includePoint(x, y, outsideDrawPoints) {
				yTo = y
			} else if isLine && !includePoint(x, y, outsideDrawPoints) {
				linePoints = append(linePoints, LinePoint{yFrom, yTo})
				isLine = false
			}
		}

		for y := 0; y <= int(height); y++ {
			color := bgcolor
			switch len(linePoints) {
			case 0:
				drawPoint = append(drawPoint, DrawPoint{x, y, bgcolor})
			case 1:
				if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = linecolor
				}
				drawPoint = append(drawPoint, DrawPoint{x, y, color})
			case 2:
				// 星の頂点にかかるときは2つのラインの間は背景色
				if x == int(math.Round(pOutside[3].X)) || x == int(math.Round(pOutside[3].X))+1 {
					if y < linePoints[0].FromY {
						color = bgcolor
					} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
						color = linecolor
					} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
						color = bgcolor
					} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
						color = linecolor
					} else {
						color = bgcolor
					}
				} else {
					if y < linePoints[0].FromY {
						color = bgcolor
					} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
						color = linecolor
					} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
						color = starColor
					} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
						color = linecolor
					} else {
						color = bgcolor
					}
				}
				drawPoint = append(drawPoint, DrawPoint{x, y, color})
			case 3:
				if y < linePoints[0].FromY {
					color = bgcolor
				} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = linecolor
				} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
					color = starColor
				} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
					color = linecolor
				} else if y > linePoints[1].ToY && y < linePoints[2].FromY {
					color = starColor
				} else if y >= linePoints[2].FromY && y <= linePoints[2].ToY {
					color = linecolor
				} else {
					color = bgcolor
				}
				drawPoint = append(drawPoint, DrawPoint{x, y, color})
			case 4:
				if y < linePoints[0].FromY {
					color = bgcolor
				} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = linecolor
				} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
					color = starColor
				} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
					color = linecolor
				} else if y > linePoints[1].ToY && y < linePoints[2].FromY {
					color = bgcolor
				} else if y >= linePoints[2].FromY && y <= linePoints[2].ToY {
					color = linecolor
				} else if y > linePoints[2].ToY && y < linePoints[3].FromY {
					color = starColor
				} else if y >= linePoints[3].FromY && y <= linePoints[3].ToY {
					color = linecolor
				} else {
					color = bgcolor
				}
				drawPoint = append(drawPoint, DrawPoint{x, y, color})
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
