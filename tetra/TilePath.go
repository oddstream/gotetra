// Copyright ©️ 2021 oddstream.games

package tetra

import "math/rand"

// TilePath describes a path of Tiles through the Grid
type TilePath struct {
	start *Tile
}

func newdir(dir uint) uint {
	r := rand.Float64()
	if r < 0.25 {
		return shiftBits(dir)
	} else if r > 0.75 {
		return unshiftBits(dir)
	}
	return dir
	// dirs := [4]uint{NORTH, EAST, SOUTH, WEST}
	// return dirs[rand.Intn(len(dirs))]
}

// Run creates a linked path of shapes on tiles
func (tp *TilePath) Run(dir uint) {
	var pos, next *Tile
	pos = tp.start
	next = dir2tile(pos, dir)
	for next != nil {
		pos.coins |= dir
		next.coins |= oppdir(dir)

		dir = newdir(dir)

		pos = next
		next = dir2tile(pos, dir)
	}
}
