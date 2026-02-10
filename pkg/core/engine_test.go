package core

import (
	"testing"

	"github.com/hkupty/mirkwood/pkg/command"
	"github.com/hkupty/mirkwood/pkg/maze"
)

// internal helper to advance through the steps.
// Returns the last processed state to aid debugging
func process(state State, actions []any) (State, error) {
	var stateCursor State
	var err error
	stateCursor = state
	for _, action := range actions {
		stateCursor, err = Step(stateCursor, action)
		if err != nil {
			return stateCursor, err
		}
		if stateCursor.IsAtFinish() {
			return stateCursor, nil
		}
	}

	if !stateCursor.IsAtFinish() {
		return stateCursor, ErrIncompletePath
	}

	return stateCursor, nil
}

func stateFromBoard(invariant levelInvariants) State {
	return State{
		invariants: invariant,
	}
}

func TestEngine(t *testing.T) {
	base := stateFromBoard(levelInvariants{
		Walls:          maze.BitBoard(0xFF2DA5B5B195C5FD),
		FinishingPoint: maze.BitBoard(1 << 55),
	})

	base.Position = 1 << 1

	final, err := process(base, stringToCommandList("ssseenneesesssse"))

	if err != nil {
		t.Error(err)
	}
	if !final.IsAtFinish() {
		t.Fatalf("\nPOS    %b\nFINISH %b\n\n Not at the end", final.Position, final.invariants.FinishingPoint)
	}

	if err = final.IsValid(); err != nil {
		t.Error(err)
	}
}

func stringToCommandList(str string) []any {
	// NOTE: This is an internal helper for testing, motivated by the fact
	// fuzz testing can only take primitive types.
	commands := make([]any, 0, len(str))
	for _, char := range str {
		switch char {
		case 's':
			commands = append(commands, command.Walk{Dir: command.South})
		case 'n':
			commands = append(commands, command.Walk{Dir: command.North})
		case 'e':
			commands = append(commands, command.Walk{Dir: command.East})
		case 'w':
			commands = append(commands, command.Walk{Dir: command.West})
		case 'm':
			commands = append(commands, command.Mark{})
		}
	}

	return commands
}

func FuzzEngine(f *testing.F) {
	base := stateFromBoard(levelInvariants{
		Walls:          maze.BitBoard(0xFF2DA5B5B195C5FD),
		FinishingPoint: maze.BitBoard(1 << 55),
	})

	base.Position = 1 << 1

	f.Add("s")
	f.Add("n")
	f.Add("w")
	f.Add("e")
	f.Add("m")
	f.Add("se")
	f.Add("nw")
	f.Add("ws")
	f.Add("em")
	f.Add("men")
	f.Fuzz(func(t *testing.T, a string) {
		state, _ := process(base, stringToCommandList(a))
		// TODO: Further verify that errors match errored states.

		if err := state.IsValid(); err != nil {
			t.Fatal(err)
		}
	})
}
