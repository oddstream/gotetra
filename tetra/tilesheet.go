// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"
	"log"

	"github.com/fogleman/gg"
)

/***
func makeTilesheet(tileSize, shapeSize int) image.Image {

	if !(tileSize > shapeSize) {
		log.Fatal("tile must be bigger than shape")
	}

	lineWidth := float64(shapeSize / 6)
	circleRadius := float64(shapeSize / 5)

	dc := gg.NewContext(tileSize*3, tileSize*2)

	// 2 coin straight
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)
	dc.MoveTo(450, 200)
	dc.LineTo(750, 200)
	dc.Stroke()

	// 2 coin L
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)
	dc.MoveTo(850, 200)
	dc.LineTo(1000, 200)
	dc.LineTo(1000, 50)
	dc.Stroke()

	// 3 coin
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)
	dc.DrawLine(450, 600, 750, 600)
	dc.DrawLine(600, 600, 600, 450)
	dc.Stroke()

	// 4 coin
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)
	dc.DrawLine(850, 600, 1150, 600)
	dc.DrawLine(1000, 750, 1000, 450)
	dc.Stroke()

	// 1 coin
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)
	dc.DrawLine(200, 750, 200, 600)
	dc.Stroke()

	dc.SetLineWidth(lineWidth)
	dc.DrawCircle(200, 600, circleRadius)
	dc.Stroke()

	return dc.Image()
}
***/

// return an image.Image that is bigger than the tile size requested so endcaps are visible
// func makeTile(coins uint, tileSize int) image.Image {

// 	tileSizeEx := tileSize + (tileSize / 6) // same as linewidth

// 	margin := float64(tileSizeEx-tileSize) / 2
// 	center := float64(tileSizeEx / 2)
// 	lineWidth := float64(tileSize / 6)
// 	circleRadius := float64(tileSize / 5)

// 	nx, ny := center, margin
// 	ex, ey := float64(tileSizeEx)-margin, center
// 	sx, sy := center, float64(tileSizeEx)-margin
// 	wx, wy := margin, center

// 	dc := gg.NewContext(tileSizeEx, tileSizeEx)
// 	dc.SetRGB(0, 0, 0)
// 	dc.SetLineWidth(lineWidth)
// 	dc.SetLineCap(gg.LineCapRound)

// 	switch coins {
// 	case 0:
// 		// explicitly do nothing
// 	case NORTH:
// 		dc.DrawLine(nx, ny, center, center-circleRadius)
// 		dc.DrawCircle(center, center, circleRadius)
// 	case EAST:
// 		dc.DrawLine(ex, ey, center+circleRadius, center)
// 		dc.DrawCircle(center, center, circleRadius)
// 	case SOUTH:
// 		dc.DrawLine(sx, sy, center, center+circleRadius)
// 		dc.DrawCircle(center, center, circleRadius)
// 	case WEST:
// 		dc.DrawLine(wx, wy, center-circleRadius, center)
// 		dc.DrawCircle(center, center, circleRadius)

// 	case NORTH | SOUTH:
// 		dc.DrawLine(nx, ny, sx, sy)
// 	case EAST | WEST:
// 		dc.DrawLine(wx, wy, ex, ey)

// 	case NORTH | EAST:
// 		dc.MoveTo(nx, ny)
// 		dc.LineTo(center, center)
// 		dc.LineTo(ex, ey)
// 	case EAST | SOUTH:
// 		dc.MoveTo(ex, ey)
// 		dc.LineTo(center, center)
// 		dc.LineTo(sx, sy)
// 	case SOUTH | WEST:
// 		dc.MoveTo(sx, sy)
// 		dc.LineTo(center, center)
// 		dc.LineTo(wx, wy)
// 	case WEST | NORTH:
// 		dc.MoveTo(wx, wy)
// 		dc.LineTo(center, center)
// 		dc.LineTo(nx, ny)
// 	case NORTH | EAST | SOUTH:
// 		dc.DrawLine(nx, ny, sx, sy)
// 		dc.DrawLine(center, center, ex, ey)
// 	case EAST | SOUTH | WEST:
// 		dc.DrawLine(wx, wy, ex, ey)
// 		dc.DrawLine(center, center, sx, sy)
// 	case SOUTH | WEST | NORTH:
// 		dc.DrawLine(nx, ny, sx, sy)
// 		dc.DrawLine(center, center, wx, wy)
// 	case WEST | NORTH | EAST:
// 		dc.DrawLine(wx, wy, ex, ey)
// 		dc.DrawLine(center, center, nx, ny)

// 	case NORTH | EAST | SOUTH | WEST:
// 		dc.DrawLine(nx, ny, sx, sy)
// 		dc.DrawLine(wx, wy, ex, ey)

// 	default:
// 		log.Fatal("makeTile called with wrong bits", coins)
// 	}
// 	dc.Stroke()

// 	return dc.Image()
// }

func makeTileCurvy(coins uint, tileSize int) image.Image {

	tileSizeEx := tileSize + (tileSize / 6) // same as linewidth

	margin := float64(tileSizeEx-tileSize) / 2
	center := float64(tileSizeEx / 2)
	lineWidth := float64(tileSize / 6)
	circleRadius := float64(tileSize / 5)

	nx, ny := center, margin
	ex, ey := float64(tileSizeEx)-margin, center
	sx, sy := center, float64(tileSizeEx)-margin
	wx, wy := margin, center

	dc := gg.NewContext(tileSizeEx, tileSizeEx)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	switch coins {
	case 0:
		// explicitly do nothing
	case NORTH:
		dc.DrawLine(nx, ny, center, center-circleRadius)
		dc.DrawCircle(center, center, circleRadius)
	case EAST:
		dc.DrawLine(ex, ey, center+circleRadius, center)
		dc.DrawCircle(center, center, circleRadius)
	case SOUTH:
		dc.DrawLine(sx, sy, center, center+circleRadius)
		dc.DrawCircle(center, center, circleRadius)
	case WEST:
		dc.DrawLine(wx, wy, center-circleRadius, center)
		dc.DrawCircle(center, center, circleRadius)

	case NORTH | SOUTH:
		dc.DrawLine(nx, ny, sx, sy)
	case EAST | WEST:
		dc.DrawLine(wx, wy, ex, ey)

	case NORTH | EAST:
		dc.MoveTo(nx, ny)
		dc.QuadraticTo(center, center, ex, ey)
	case EAST | SOUTH:
		dc.MoveTo(ex, ey)
		dc.QuadraticTo(center, center, sx, sy)
	case SOUTH | WEST:
		dc.MoveTo(sx, sy)
		dc.QuadraticTo(center, center, wx, wy)
	case WEST | NORTH:
		dc.MoveTo(wx, wy)
		dc.QuadraticTo(center, center, nx, ny)

	case NORTH | EAST | SOUTH:
		dc.MoveTo(nx, ny)
		dc.QuadraticTo(center, center, ex, ey)
		dc.QuadraticTo(center, center, sx, sy)
	case EAST | SOUTH | WEST:
		dc.MoveTo(ex, ey)
		dc.QuadraticTo(center, center, sx, sy)
		dc.QuadraticTo(center, center, wx, wy)
	case SOUTH | WEST | NORTH:
		dc.MoveTo(sx, sy)
		dc.QuadraticTo(center, center, wx, wy)
		dc.QuadraticTo(center, center, nx, ny)
	case WEST | NORTH | EAST:
		dc.MoveTo(wx, wy)
		dc.QuadraticTo(center, center, nx, ny)
		dc.QuadraticTo(center, center, ex, ey)

	case NORTH | EAST | SOUTH | WEST:
		dc.MoveTo(nx, ny)
		dc.QuadraticTo(center, center, ex, ey)
		dc.QuadraticTo(center, center, sx, sy)
		dc.QuadraticTo(center, center, wx, wy)
		dc.QuadraticTo(center, center, nx, ny)

	default:
		log.Fatal("makeTile called with wrong bits", coins)
	}
	dc.Stroke()

	return dc.Image()
}
