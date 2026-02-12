package mazeview

import (
	"testing"

	"github.com/hkupty/mirkwood/pkg/maze"
)

func TestCanWriteMaze(t *testing.T) {
	maze := New(maze.SampleBlueprint)
	str := maze.View()
	if len(str) > 0 {
		t.Fatal(len(str))
	}
}
