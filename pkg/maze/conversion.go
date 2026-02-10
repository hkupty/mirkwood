package maze

// GridToBitBoard converts a MazeGrid to a BitBoard representation.
// Returns the bitboard where each bit represents a wall (1) or path (0).
func GridToBitBoard(grid MazeGrid) BitBoard {
	var bitBoard uint64

	// NOTE: Very important here that rows and columns are expected to always be 8x8
	for jx, row := range grid {
		for ix, col := range row {
			if col {
				bitBoard |= 1 << (8*jx + ix)
			}
		}
	}

	return BitBoard(bitBoard)
}

// PosToBit converts grid coordinates to a bit position (0-63).
func PosToBit(row, col uint8) uint8 {
	return row*8 + col
}

// BitToPos converts a bit position to grid coordinates.
func BitToPos(bit uint8) (row, col uint8) {
	return bit / 8, bit % 8
}
