package core

import (
	"fmt"

	"github.com/hkupty/mirkwood/pkg/command"
)

func Step(state State, action any) (State, error) {
	var nextState State
	var err error

	switch v := action.(type) {
	case command.Walk:
		nextState, err = state.Move(v.Dir)
	case command.Mark:
		nextState = state.ToggleMark()
	default:
		return state, fmt.Errorf("unknown action type: %T", action)
	}

	if err != nil {
		return state, err
	}

	return nextState, nil
}
