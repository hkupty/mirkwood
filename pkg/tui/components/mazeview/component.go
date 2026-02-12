package mazeview

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	model := Model{
		state:  core.NewStateFromBlueprint(bp),
		buffer: composite.NewBuffer(12, 5),
	}
	state, _ := model.Update(nil)
	return state
}

var trees = []string{
	"█", "▓", "▒", "░",
}

// TODO: Transitions and ornaments
var chars = []string{
	" ",
	"█",
	"●",
	"",
	"░",
	"",
	"✪",
	"",
}

var bgArray = []lipgloss.Color{
	styles.PathBg,
	styles.PathBg,
	styles.VisitedPathBg,
	styles.VisitedPathBg,
}

var fgArray = []lipgloss.Color{
	styles.PlayerFg,
	styles.MarkFg,
	styles.PlayerFg,
	styles.VisitedMarkFg,
}

func (m Model) Update(action any) (Model, error) {
	if action != nil {
		state, err := core.Step(m.state, action)
		if err != nil {
			return m, err
		}
		m.state = state
	}
	m.buffer.Composite(m.state)
	return m, nil
}

func (m Model) View() string {
	var buffer strings.Builder
	style := lipgloss.NewStyle()

	for _, row := range m.buffer.Cells {
		for _, cell := range row {
			identity := uint(cell) & 0b1111
			if identity == 1 {
				decor := cell >> 4
				shade := decor & 0b11
				buffer.WriteString(style.Background(styles.WallBg).Foreground(styles.WallFg).Render(trees[shade]))
			} else {
				ix := (identity >> 2)
				buffer.WriteString(style.Background(bgArray[ix]).Foreground(fgArray[ix]).Render(chars[identity&0b111]))
			}

		}
		buffer.WriteRune('\n')
	}
	return buffer.String()
}
