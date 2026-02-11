package composite

type Buffer struct {
	XRes  int
	YRes  int
	Cells [][]Cell
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
//	└─┴─┴─┴────────╴Unused (yet)
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
	cells := make([][]Cell, col_size)
	rows := make([]Cell, row_size*col_size)

	for ix := range col_size {
		cells[ix] = rows[ix*row_size : (ix+1)*row_size]
		for jx, _ := range cells[ix] {
			cells[ix][jx] = NewCell(0, true, true)
		}
	}

	return Buffer{
		XRes:  xres,
		YRes:  yres,
		Cells: cells,
	}
}
