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

type Star struct {
	// 画像サイズ
	width  float64
	height float64

	// 画像を保存するファイル名
	filename string

	// 中心のx,y座標
	centerX float64
	centerY float64

	// 背景色とラインカラー
	bgColor   color.Color
	lineColor color.Color
	starColor color.Color
}

func newStar(width, height float64, filename string, bgColor, lineColor, starColor color.Color) *Star {
	return &Star{
		width:     width,
		height:    height,
		filename:  filename,
		centerX:   width / 2,
		centerY:   height / 2,
		bgColor:   bgColor,
		lineColor: lineColor,
		starColor: starColor,
	}
}

func (s *Star) drawStar() {
	m := image.NewRGBA(image.Rect(0, 0, int(s.width), int(s.height)))
	draw.Draw(m, m.Bounds(), &image.Uniform{s.bgColor}, image.Point{}, draw.Src)

	pOutside := s.pentagonOutside()
	pInside := s.pentagonInside()

	// 星の外枠の座標を求める
	var starOutsidePoints []*Point
	for i := 0; i < len(pOutside); i++ {
		starOutsidePoints = append(starOutsidePoints, s.getCoordinates(pOutside[i], pInside[i])...)

		if i+1 == len(pOutside) {
			// pointsInsideの最後の点はpointsOutsideの最初の点につなげる
			starOutsidePoints = append(starOutsidePoints, s.getCoordinates(pInside[i], pOutside[0])...)
		} else {
			starOutsidePoints = append(starOutsidePoints, s.getCoordinates(pInside[i], pOutside[i+1])...)
		}
	}

	// 星の外枠の座標
	outsideDrawPoints := s.convertPint2DrawPoint(starOutsidePoints, s.lineColor)

	// 各ますに対する色を決定する
	drawPoints := s.registerColor(outsideDrawPoints, pOutside)

	// 実際に点を打ってファイルに書き出す
	for _, p := range drawPoints {
		m.Set(p.X, p.Y, p.Color)
	}
	file, err := os.Create(s.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, m)
	if err != nil {
		log.Fatal(err)
	}
}

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

// 星の外側の頂点座標を求める
func (s *Star) pentagonOutside() []*Point {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := s.width / 2.0

	// 中心から頂点への直線同士が成す角度
	radian := math.Pi * 2 / vertexNum

	var points []*Point
	for i := 0.0; i < vertexNum; i++ {
		x := s.centerX + R*math.Cos(radian*i)
		y := s.centerY + R*math.Sin(radian*i)
		points = append(points, &Point{x, y})
		fmt.Printf("外側：(x,y)=(%d,%d)\n", int(math.Round(x)), int(math.Round(y)))
	}
	return points
}

// 星の内側の頂点座標を求める
func (s *Star) pentagonInside() []*Point {
	// 頂点の数
	vertexNum := 5.0

	// 中心から頂点までの距離
	R := s.width / 4.0

	// 外側の1/2の角度を求める
	radian := math.Pi / vertexNum

	var points []*Point
	for i := 1.0; i < vertexNum*2; i = i + 2 { // 外側と角度をずらしつつ幅を揃える
		x := s.centerX + R*math.Cos(radian*i)
		y := s.centerY + R*math.Sin(radian*i)
		points = append(points, &Point{x, y})
		fmt.Printf("内側：(x,y)=(%d,%d)\n", int(math.Round(x)), int(math.Round(y)))
	}

	return points
}

// 2点を受け取ってその間の座標を返す
func (s *Star) getCoordinates(xy1, xy2 *Point) []*Point {
	dx := xy2.X - xy1.X
	dy := xy2.Y - xy1.Y

	// 距離
	length := math.Sqrt(dx*dx + dy*dy)

	// ラジアン
	radian := math.Atan2(dy, dx)

	// 2点間の座標を格納するスライス
	var points []*Point

	for l := 0.0; l < length; l++ {
		// x座標とy座標を計算
		x := xy1.X + l*math.Cos(radian)
		y := xy1.Y + l*math.Sin(radian)

		// ビットマップ外の点は描写しない
		if (x >= 0 && x < s.width) && (y >= 0 && y < s.height) {
			points = append(points, &Point{x, y})
		}
	}

	return points
}

// float64で定義されている座標(Point型)の一覧とその座標群の色を受け取り、DrawPoint型に詰め直して返す
// 同時に重複も削除する
func (s *Star) convertPint2DrawPoint(points []*Point, color color.Color) []*DrawPoint {
	var drawPoints []*DrawPoint
	for _, p := range points {
		drawPoints = append(drawPoints, &DrawPoint{int(math.Round(p.X)), int(math.Round(p.Y)), color})
	}
	return s.deleteDuplicate(drawPoints)
}

// Sliceの重複を削除、x軸→y軸の値の昇順にソートする
func (s *Star) deleteDuplicate(old []*DrawPoint) []*DrawPoint {
	m := make(map[DrawPoint]bool)
	var new []*DrawPoint

	for _, point := range old {
		if !m[*point] {
			m[*point] = true
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

// 各マス目の色を決める
func (s *Star) registerColor(outsideDrawPoints []*DrawPoint, pOutside []*Point) []*DrawPoint {
	var drawPoints []*DrawPoint

	for x := 0; x <= int(s.width); x++ {
		// 直前のマスがライン上かどうか保持する変数
		isLine := false
		yTo := 0
		yFrom := 0
		var linePoints []*LinePoint

		// 星の枠線のy座標を求める
		for y := 0; y <= int(s.height); y++ {
			if !isLine && s.includePoint(x, y, outsideDrawPoints) {
				// ラインじゃないところからライン上に変わったらFromToを両方設定する
				yFrom = y
				yTo = y
				isLine = true
			} else if isLine && s.includePoint(x, y, outsideDrawPoints) {
				// 前のマスに引き続きライン上なら、Toの値を更新する
				yTo = y
			} else if isLine && !s.includePoint(x, y, outsideDrawPoints) {
				// ラインから外れたらライン情報を格納する
				linePoints = append(linePoints, &LinePoint{yFrom, yTo})
				isLine = false
			}
		}

		for y := 0; y <= int(s.height); y++ {
			color := s.bgColor
			switch len(linePoints) {
			case 0:
				drawPoints = append(drawPoints, &DrawPoint{x, y, s.bgColor})
			case 1:
				// ラインが1箇所しかないときは、ライン以外は背景色
				if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = s.lineColor
				}
				drawPoints = append(drawPoints, &DrawPoint{x, y, color})
			case 2:
				// ラインが2箇所有るときは、背景→ライン→星→ライン→背景のパターンと、背景→頂点(ライン)→背景→頂点(ライン)→背景のパターンがある
				if x == int(math.Round(pOutside[3].X)) || x == int(math.Round(pOutside[3].X))+1 {
					// 背景→頂点(ライン)→背景→頂点(ライン)→背景のパターン
					if y < linePoints[0].FromY {
						color = s.bgColor
					} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
						color = s.lineColor
					} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
						color = s.bgColor
					} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
						color = s.lineColor
					} else {
						color = s.bgColor
					}
				} else {
					// 背景→ライン→星→ライン→背景のパターン
					if y < linePoints[0].FromY {
						color = s.bgColor
					} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
						color = s.lineColor
					} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
						color = s.starColor
					} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
						color = s.lineColor
					} else {
						color = s.bgColor
					}
				}
				drawPoints = append(drawPoints, &DrawPoint{x, y, color})
			case 3:
				if y < linePoints[0].FromY {
					color = s.bgColor
				} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = s.lineColor
				} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
					color = s.starColor
				} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
					color = s.lineColor
				} else if y > linePoints[1].ToY && y < linePoints[2].FromY {
					color = s.starColor
				} else if y >= linePoints[2].FromY && y <= linePoints[2].ToY {
					color = s.lineColor
				} else {
					color = s.bgColor
				}
				drawPoints = append(drawPoints, &DrawPoint{x, y, color})
			case 4:
				if y < linePoints[0].FromY {
					color = s.bgColor
				} else if y >= linePoints[0].FromY && y <= linePoints[0].ToY {
					color = s.lineColor
				} else if y > linePoints[0].ToY && y < linePoints[1].FromY {
					color = s.starColor
				} else if y >= linePoints[1].FromY && y <= linePoints[1].ToY {
					color = s.lineColor
				} else if y > linePoints[1].ToY && y < linePoints[2].FromY {
					color = s.bgColor
				} else if y >= linePoints[2].FromY && y <= linePoints[2].ToY {
					color = s.lineColor
				} else if y > linePoints[2].ToY && y < linePoints[3].FromY {
					color = s.starColor
				} else if y >= linePoints[3].FromY && y <= linePoints[3].ToY {
					color = s.lineColor
				} else {
					color = s.bgColor
				}
				drawPoints = append(drawPoints, &DrawPoint{x, y, color})
			}
		}
	}

	return drawPoints
}

// (x,y)が[]*DrawPointに含まれているかどうかを判定する
func (s *Star) includePoint(x, y int, points []*DrawPoint) bool {
	for _, p := range points {
		if x == p.X && y == p.Y {
			return true
		}
	}
	return false
}
