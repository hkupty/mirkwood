package core

import (
	"fmt"

	"github.com/hkupty/mirkwood/pkg/command"
)

func Step(state State, action any) (State, error) {
	switch v := action.(type) {
	case command.Walk:
		return state.Move(v.Dir)
	case command.Mark:
		return state.ToggleMark(), nil
	}

	return state, fmt.Errorf("unknown action type: %T", action)
}

func Process(state State, actions []any) (State, error) {

	// NOTE: This is mostly for internal processing. While we don't want to re-implement this in TUI,
	// we need to have a proper step-refresh screen synchrony, with a possible time delay (200~800ms)
	// so the result is visible.

	var stateCursor State
	var err error
	stateCursor = state
	for _, action := range actions {
		stateCursor, err = Step(stateCursor, action)
		if err != nil {
			return state, err
		}
		if stateCursor.IsAtFinish() {
			return stateCursor, nil
		}
	}

	if !stateCursor.IsAtFinish() {
		return state, ErrIncompletePath
	}

	return stateCursor, nil

}
