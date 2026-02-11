package core

import "github.com/hkupty/mirkwood/pkg/maze"

// State represents the runtime state of a level.
// It tracks the player's position, visited cells, marks, and step count.
type State struct {
	// Position is a single-bit bitboard indicating where the player is
	Position maze.BitBoard

	// VisitedPath tracks all cells the player has stepped on
	VisitedPath maze.BitBoard

	// Marks tracks which cells the player has marked
	Marks maze.BitBoard

	// StepsCounter tracks how many moves the player has made
	StepsCounter uint16

	// Static level data (walls, finish point)
	Invariants LevelInvariants
}

// LevelInvariants holds static level data that doesn't change during play
type LevelInvariants struct {
	// Walls bitboard (static, immutable)
	Walls maze.BitBoard

	// FinishingPoint is the target position to reach
	FinishingPoint maze.BitBoard
}
