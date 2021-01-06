// Copyright ©️ 2020 oddstream.games

package tetra

import (
	"image"
)

// InRect returns true of px,py is within Rect returned by function parameter
func InRect(pt image.Point, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return pt.X > x0 && pt.Y > y0 && pt.X < x1 && pt.Y < y1
}
