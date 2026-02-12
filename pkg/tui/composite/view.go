package composite

type Buffer struct {
	XRes    int
	YRes    int
	Context [][]ContextCell
	Dirty   uint64
	Cells   [][]Cell
}

type CellIdentity uint8

const (
	Path CellIdentity = iota
	Wall
	Player
	reserved
)

// a cell is a bit flag, in which each bit position represents some information
// and the combination of bits provides the rendering layer with enough information to
// display the right information. Note that a cell can never be Player and Wall at the same time,
// so the values 3, 7, 11 and 15 are invalid.
//
//	   ┌───────────╴Cell Decoration
//	   │       ┌───╴Block Identity (All cells of the same block will have this the same)
//	┌─┬┴┬─┐ ┌─┬┴┬─┐
//	0 0 0 0 0 0 0 0
//	│ │ │ │ │ │ └┬┘
//	│ │ │ │ │ │  └─╴Path=00		Wall=01			Player=10		Reserved=11
//	│ │ │ │ │ └────╴Marked=1	Unmarked=0
//	│ │ │ │ └──────╴Visited=1	Unvisited=0
//	└─┴─┴─┴────────╴Accent (identity-dependent)
//
// When a cell represents a wall, its identity will be always 0001
// When a cell represents a path, its identity can be:
//   - 0000 (unvisited, unmarked path)
//   - 0100 (unvisited, marked path)
//   - 1010 (visited, player-standing path)
//   - 1110 (visited, marked, player-standing path)
//   - 1100 (visited, marked path)
//
// this means that identity >> 2 gives us [00, 01, 10, 11] for using as decoration index
// Player is only present on visited paths always.
type Cell uint8

func NewCell(identity CellIdentity, marked bool, visited bool) Cell {
	cell := Cell(identity)

	if marked {
		cell |= 1 << 2
	}

	if visited {
		cell |= 1 << 3
	}

	return cell
}

func NewBuffer(xres, yres int) Buffer {
	row_size := xres * 8
	col_size := yres * 8
	context := make([][]ContextCell, 8)
	context_rows := make([]ContextCell, 64)

	for ix := range 8 {
		context[ix] = context_rows[ix*8 : (ix+1)*8]
	}

	cells := make([][]Cell, col_size)
	rows := make([]Cell, row_size*col_size)

	for ix := range col_size {
		cells[ix] = rows[ix*row_size : (ix+1)*row_size]
	}

	return Buffer{
		XRes:    xres,
		YRes:    yres,
		Context: context,
		Cells:   cells,
	}
}
