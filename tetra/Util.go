// Copyright ©️ 2020 oddstream.games

package tetra

// InRect returns true of px,py is within Rect returned by function parameter
func InRect(px int, py int, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return px > x0 && py > y0 && px < x1 && py < y1
}
