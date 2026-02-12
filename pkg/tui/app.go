package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hkupty/mirkwood/pkg/command"
	"github.com/hkupty/mirkwood/pkg/maze"
	"github.com/hkupty/mirkwood/pkg/tui/components/mazeview"
)

type model struct {
	maze mazeview.Model
}

func initialModel() model {
	return model{
		maze: mazeview.New(maze.SampleBlueprint),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "j":
			newM, err := m.maze.Update(command.Walk{Dir: command.South})
			if err != nil {
				return m, nil
			}

			m.maze = newM

		case "k":
			newM, err := m.maze.Update(command.Walk{Dir: command.North})
			if err != nil {
				return m, nil
			}

			m.maze = newM

		case "l":
			newM, err := m.maze.Update(command.Walk{Dir: command.East})
			if err != nil {
				return m, nil
			}

			m.maze = newM

		case "h":
			newM, err := m.maze.Update(command.Walk{Dir: command.West})
			if err != nil {
				return m, nil
			}

			m.maze = newM

		case "m":
			newM, err := m.maze.Update(command.Mark{})
			if err != nil {
				return m, nil
			}

			m.maze = newM
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
		return m, nil
	}

	return m, nil
}

func (m model) View() string {

	return m.maze.View()
}

func MainLoop() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
