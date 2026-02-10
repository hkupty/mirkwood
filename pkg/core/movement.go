package core

import (
	"errors"
	"math/bits"

	"github.com/hkupty/mirkwood/pkg/command"
	"github.com/hkupty/mirkwood/pkg/maze"
)

var (
	// ErrInvalidState indicates the state has become corrupted (e.g., multiple positions)
	ErrInvalidState = errors.New("invalid game state")

	// ErrHitWall indicates the player attempted to move into a wall
	ErrHitWall = errors.New("you hit a tree")

	// ErrStepLimit indicates the player exceeded the maximum step count
	ErrStepLimit = errors.New("step limit exceeded")

	// ErrAlreadyComplete indicates the level is already finished
	ErrAlreadyComplete = errors.New("level already complete")
)

// Move attempts to move the player in the given direction.
// Returns the new state and an error if the move is invalid.
func (s State) Move(dir command.Direction) (State, error) {
	// Verify we have exactly one position bit set
	if bits.OnesCount64(uint64(s.Position)) != 1 {
		return s, ErrInvalidState
	}

	// Calculate next position based on direction
	var nextPos maze.BitBoard
	switch dir {
	case command.North:
		nextPos = s.Position >> 8
	case command.South:
		nextPos = s.Position << 8
	case command.East:
		nextPos = s.Position << 1
	case command.West:
		nextPos = s.Position >> 1
	}

	// Check for wall collision
	if nextPos&s.invariants.Walls != 0 {
		return s, ErrHitWall
	}

	// Create new state (immutable update)
	return State{
		Position:     nextPos,
		VisitedPath:  s.VisitedPath | nextPos,
		Marks:        s.Marks,
		StepsCounter: s.StepsCounter + 1,
		invariants:   s.invariants,
	}, nil
}

// ToggleMark toggles a mark on the current position.
// Returns the new state with the mark toggled.
func (s State) ToggleMark() State {
	return State{
		Position:     s.Position,
		VisitedPath:  s.VisitedPath,
		Marks:        s.Marks ^ s.Position,
		StepsCounter: s.StepsCounter,
		invariants:   s.invariants,
	}
}

// IsAtFinish returns true if the player is at the finishing point
func (s State) IsAtFinish() bool {
	return s.Position&s.invariants.FinishingPoint != 0
}

// MarkCount returns the number of marked cells
func (s State) MarkCount() int {
	return bits.OnesCount64(uint64(s.Marks))
}

// IsValid validates that the state satisfies all invariants
func (s State) IsValid() error {
	// Check exactly one position
	if bits.OnesCount64(uint64(s.Position)) != 1 {
		return ErrInvalidState
	}

	// Check no overlap between marks and walls
	if s.Marks&s.invariants.Walls != 0 {
		return ErrInvalidState
	}

	// Check no overlap between visited path and walls
	if s.VisitedPath&s.invariants.Walls != 0 {
		return ErrInvalidState
	}

	// Check position is not on a wall
	if s.Position&s.invariants.Walls != 0 {
		return ErrInvalidState
	}

	return nil
}
