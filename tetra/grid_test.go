// Copyright ©️ 2020 oddstream.games

package tetra

// load png decoder in main package
import (
	_ "image/png"
	"testing"
)

func TestTileLinking(t *testing.T) {
	g := NewGrid("puzzle", 0, 0)
	t00 := g.findTile(0, 0)
	if t00 == nil {
		t.Error("cannot find t00")
	}
	t10 := g.findTile(1, 0)
	if t10 == nil {
		t.Error("cannot find t10")
	}
	if t00.edges[1] != t10 {
		t.Error("t00 t10 not linked")
	}
	if t10.edges[3] != t00 {
		t.Error("t10 t0 not linked")
	}
	t01 := g.findTile(0, 1)
	if t01 == nil {
		t.Error("cannot find t01")
	}
	if t00.edges[2] != t01 {
		t.Error("t00 t01 not linked")
	}
	if t00.edges[2].edges[0] != t00 {
		t.Error("t00.S.N not linked")
	}
}
