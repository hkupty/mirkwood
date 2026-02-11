package core

import (
	"github.com/hkupty/mirkwood/pkg/maze"
)

// NewStateFromBlueprint creates a new game State from a LevelBlueprint.
// This is the primary way to initialize a level's runtime state.
func NewStateFromBlueprint(bp maze.LevelBlueprint) State {
	walls := maze.GridToBitBoard(bp.Grid)
	startPos := maze.BitBoard(1 << bp.StartingPoint)
	finishPos := maze.BitBoard(1 << bp.FinishingPoint)

	return State{
		Position:     startPos,
		VisitedPath:  startPos,
		Marks:        0,
		StepsCounter: 0,
		Invariants: LevelInvariants{
			Walls:          walls,
			FinishingPoint: finishPos,
		},
	}
}
