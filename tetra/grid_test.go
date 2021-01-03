// Copyright ©️ 2020 oddstream.games

package tetra

// load png decoder in main package
import (
	_ "image/png"
	"testing"
)

func TestTileLinking(t *testing.T) {
	g := NewGrid("puzzle")
	t00 := g.findTile(0, 0)
	if t00 == nil {
		t.Error("cannot find t00")
	}
	t10 := g.findTile(1, 0)
	if t10 == nil {
		t.Error("cannot find t10")
	}
	if t00.E != t10 {
		t.Error("t00 t10 not linked")
	}
	if t10.W != t00 {
		t.Error("t10 t0 not linked")
	}
	t01 := g.findTile(0, 1)
	if t01 == nil {
		t.Error("cannot find t01")
	}
	if t00.S != t01 {
		t.Error("t00 t01 not linked")
	}
	if t00.S.N != t00 {
		t.Error("t00.S.N not linked")
	}
}
