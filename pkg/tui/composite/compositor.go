package composite

import (
	"github.com/hkupty/mirkwood/pkg/core"
)

func (buffer *Buffer) Composite(state core.State) {
	for ix := range 64 {
		var mask uint64 = 1 << ix
		wall := (uint64(state.Invariants.Walls) & mask) >> ix
		player := (uint64(state.Position) & mask) >> ix
		marks := (uint64(state.Marks)&mask)>>ix == 1
		visited := (uint64(state.VisitedPath)&mask)>>ix == 1
		cell := NewCell(CellIdentity(wall|(player<<1)), marks, visited)

		logicalX := ix % 8
		logicalY := ix / 8
		cellX := logicalX * buffer.XRes
		cellY := logicalY * buffer.YRes

		for yoff := range buffer.YRes {
			y := cellY + yoff
			for xoff := range buffer.XRes {
				x := cellX + xoff
				buffer.Cells[y][x] = cell
			}
		}
	}
}
