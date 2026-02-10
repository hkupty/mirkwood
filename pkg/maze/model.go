package maze

// BitBoard is a (max) 8x8 board linearized into a single integer for cache friendliness.
// It can represent walls (1 = wall), visited paths, or marked cells.
type BitBoard uint64

// MazeGrid is a matrix representation used for construction and analysis.
// true = wall, false = path
type MazeGrid [][]bool

// LevelBlueprint holds the static definition of a maze level.
// It includes the grid layout, start/end positions, and win condition.
type LevelBlueprint struct {
	// Key identifies the level (file number, seed, etc.)
	Key uint32

	// Grid is the 8x8 maze layout (true = wall)
	Grid MazeGrid

	// StartingPoint is the bit position (0-63) where the player begins
	StartingPoint uint8

	// FinishingPoint is the bit position (0-63) the player must reach
	FinishingPoint uint8

	// WinCondition defines what must be satisfied to complete the level
	WinCondition WinCondition
}

// WinCondition specifies how a level is completed
type WinCondition struct {
	// RequiredMarks is the number of cells that must be marked (0 = no requirement)
	RequiredMarks uint8

	// MaxSteps is the maximum allowed steps (0 = unlimited)
	MaxSteps uint16
}

// WinCondition types for convenience
var (
	// SimpleExit only requires reaching the exit
	SimpleExit = WinCondition{}
)

// A typical maze grid could look roughly like:
// [[true, false, true,  true,  true,  true,  true,  true],
//
//	[true, false, true,  false, false, false, true,  true],
//	[true, false, true,  false, true,  false, false, true],
//	[true, false, false, false, true,  true,  false, true],
//	[true, false, true,  false, true,  true,  false, true],
//	[true, false, true,  false, false, true,  false, true],
//	[true, false, true,  true,  false, true,  false, false],
//	[true, true,  true,  true,  true,  true,  true,  true]]
//
// In the MazeGrid above, the entrance is at [0][1] and the exit would be at [6][7]
// The BitBoard counterpart of that would be:
// 0b11111111_00101101_10100101_10110101_10110001_10010101_11000101_11111101
// 0xFF2DA5B5B195C5FD
// 18387535053410649597
