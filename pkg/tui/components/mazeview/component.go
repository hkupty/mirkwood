package mazeview

import (
	"strings"

	"github.com/hkupty/mirkwood/pkg/core"
	"github.com/hkupty/mirkwood/pkg/maze"
	"github.com/hkupty/mirkwood/pkg/tui/composite"
	"github.com/hkupty/mirkwood/pkg/tui/styles"
)

type Model struct {
	state  core.State
	buffer composite.Buffer
}

func New(bp maze.LevelBlueprint) Model {
	return Model{
		state:  core.NewStateFromBlueprint(bp),
		buffer: composite.NewBuffer(5, 3),
	}
}

var sprites = []string{
	styles.Path.Render(" "),
	styles.Walls.Render("█"),
	styles.Player.Render("●"),
	"",
	styles.Mark.Render(" "),
	"",
	styles.Mark.Render("●"),
	"",
	styles.Visited.Render(" "),
	"",
	styles.Visited.Render("●"),
	"",

	styles.VisitedMark.Render(" "),
	"",
	styles.VisitedMark.Render("●"),
	"",
}
const cellLineWidth = 5 * 8 * 2

func (m Model) Update() {
	m.buffer.Composite(m.state)
}

func (m Model) View() string {
	var sb strings.Builder

	for _, row := range m.buffer.Cells {
		for _, cell := range row {
			sb.WriteString(sprites[int(cell)&0b1111])
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
