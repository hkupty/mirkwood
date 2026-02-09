package maze

import (
	"errors"
	"math/bits"
)

type Direction uint8

const (
	North Direction = iota
	South
	East
	West
)

// A BitBoard is a (max) 8x8 board linearized into a single integer for cache friendliness.
// It can represent a maze (wall = 1; path = 0), a traversed path (visited = 1; unvisited = 0) or something else
// A BitBoard is mirrored from a MazeGrid.
type BitBoard uint64

// The MazeGrid is a matrix of walls and is used as a foundation to build the runtime BitBoard for the maze
type MazeGrid [][]bool

// A typical maze grid could look roughly like:
// [[true, false, true,  true,  true,  true,  true,  true],
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

type LevelState struct {
	Position     BitBoard
	VisitedPath  BitBoard
	Marks        BitBoard
	StepsCounter uint16
	invariants   LevelInvariants
}

type LevelInvariants struct {
	Walls          BitBoard
	FinishingPoint BitBoard
}

func (ls LevelState) Walk(dir Direction) (LevelState, error) {
	var nextPos BitBoard
	switch dir {
	case North:
		nextPos = ls.Position >> 8
	case South:
		nextPos = ls.Position << 8
	case East:
		nextPos = ls.Position << 1
	case West:
		nextPos = ls.Position >> 1
	}

	if bits.OnesCount64(uint64(ls.Position)) != 1 {
		return ls, errors.New("Invalid state")
	}

	// NOTE: The walls are going to wrap the edges of the map as a frame, so therefore
	// this check ensures we're not over/under flowing to a row below/above.
	if nextPos&ls.invariants.Walls != 0 {
		// TODO: Make this a global error
		return ls, errors.New("You hit a tree!")
	}

	nextState := LevelState{
		Marks:        ls.Marks,
		VisitedPath:  ls.VisitedPath | nextPos,
		Position:     nextPos,
		StepsCounter: ls.StepsCounter + 1,
		invariants:   ls.invariants,
	}

	return nextState, nil
}

func (ls LevelState) ToggleMark() (LevelState, error) {
	nextState := LevelState{
		Marks:        ls.Marks ^ ls.Position,
		VisitedPath:  ls.VisitedPath,
		Position:     ls.Position,
		StepsCounter: ls.StepsCounter,
		invariants:   ls.invariants,
	}
	return nextState, nil
}

type LevelBlueprint struct {
	// Key is a numeric value that either works as a seed or is the file key.
	// Those are decided outside of this struct - first a list of file keys are loaded
	// (01.txt, 02.txt, ... -> 1, 2, ...) into a list and the file name is turned into a key
	// to be used here or a map is generated and the seed is supplied as the key.
	// Therefore, this is for informational purposes only and setting a Key to 1
	// won't automatically load the file `01.txt`.
	Key            uint32
	Grid           MazeGrid
	StartingPoint  uint64 // maybe a [2]uint8? i.e. [0, 7]
	FinishingPoint uint64 // maybe a [2]uint8?
}

func (lb *LevelBlueprint) materialize() LevelState {
	var bitBoard uint64

	for jx, row := range lb.Grid {
		for ix, col := range row {
			if col {
				bitBoard |= 1 << (8*jx + ix)
			}
		}
	}

	return LevelState{
		Position:     BitBoard(lb.StartingPoint),
		VisitedPath:  0,
		Marks:        0,
		StepsCounter: 0,
		invariants: LevelInvariants{
			Walls:          BitBoard(bitBoard),
			FinishingPoint: BitBoard(lb.FinishingPoint),
		},
	}
}
