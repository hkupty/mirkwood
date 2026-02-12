package composite

import (
	"github.com/hkupty/mirkwood/pkg/core"
)

func (buffer *Buffer) Contextualize(state core.State) {
	for ix := range 64 {
		logicalX := ix % 8
		logicalY := ix / 8
		var mask uint64 = 1 << ix
		wall := (uint64(state.Invariants.Walls) & mask) >> ix
		player := (uint64(state.Position) & mask) >> ix
		marks := (uint64(state.Marks)&mask)>>ix == 1
		visited := (uint64(state.VisitedPath)&mask)>>ix == 1
		contextCell := NewContextCell(CellIdentity(wall|(player<<1)), marks, visited)
		if contextCell != buffer.Context[logicalY][logicalX] {
			buffer.Context[logicalY][logicalX] = contextCell
			buffer.Dirty |= mask
		}
	}
}

func (buffer *Buffer) Raster() {
	for logicalY, row := range buffer.Context {
		cellY := logicalY * buffer.YRes

		var upperRow []ContextCell
		var lowerRow []ContextCell

		if logicalY != 0 {
			upperRow = buffer.Context[logicalY-1]
		}
		if logicalY < 7 {
			lowerRow = buffer.Context[logicalY+1]
		}

		for logicalX, ctx := range row {
			cellX := logicalX * buffer.XRes
			offset := (logicalY*8 + logicalX)

			if (buffer.Dirty&(1<<offset))>>offset == 1 {

				var topNeighbor *ContextCell
				var botNeighbor *ContextCell
				var rightNeighbor *ContextCell
				var leftNeighbor *ContextCell

				if upperRow != nil {
					topNeighbor = &upperRow[logicalX]
				}

				if lowerRow != nil {
					botNeighbor = &lowerRow[logicalX]
				}

				if logicalX != 0 {
					leftNeighbor = &row[logicalX-1]
				}
				if logicalX < 7 {
					rightNeighbor = &row[logicalX+1]
				}

				for yoff := range buffer.YRes {
					y := cellY + yoff
					ydecor := 0b0000
					if topNeighbor == nil || topNeighbor.Type != ctx.Type {
						if botNeighbor != nil && botNeighbor.Type == ctx.Type {
							ydecor = min(yoff, 3)
						} else {
							ydecor = min(min(yoff, buffer.YRes-1-yoff), 3)
						}
					} else if botNeighbor == nil || botNeighbor.Type != ctx.Type {
						if topNeighbor.Type == ctx.Type {
							ydecor = min(buffer.YRes-1-yoff, 3)
						} else {
							ydecor = min(min(yoff, buffer.YRes-1-yoff), 3)
						}
					} else {
						ydecor = 3
					}

					for xoff := range buffer.XRes {
						x := cellX + xoff
						cell := ctx.ToCell()

						xdecor := 0b0000

						if leftNeighbor == nil || leftNeighbor.Type != ctx.Type {
							if rightNeighbor != nil && rightNeighbor.Type == ctx.Type {
								xdecor = min(xoff, 3)
							} else {
								xdecor = min(min(xoff, buffer.XRes-1-xoff), 3)
							}
						} else if rightNeighbor == nil || rightNeighbor.Type != ctx.Type {
							if leftNeighbor.Type == ctx.Type {
								xdecor = min(buffer.XRes-1-xoff, 3)
							} else {
								xdecor = min(min(xoff, buffer.XRes-1-xoff), 3)
							}
						} else {
							xdecor = 3
						}
						buffer.Cells[y][x] = cell | Cell(min(ydecor, xdecor)<<4)

					}
				}
			}
		}
	}
}

func (buffer *Buffer) Composite(state core.State) {
	buffer.Contextualize(state)
	buffer.Raster()
	buffer.Dirty = 0 // Clear the flags for the next frame
}
