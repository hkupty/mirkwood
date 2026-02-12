package styles

import "github.com/charmbracelet/lipgloss"

// PathChar is the character used for filling path cells.
const (
	PathChar = ' '
)

// NOTE: Marked, unless marked from start, implies visited;
// NOTE: Colors are difficult to get right. Player might be a Foreground-only modifier,
// while Path/Visited/Marked might be a background-only modifier

var (
	WallFg        = lipgloss.Color("#385831")
	WallBg        = lipgloss.Color("#1B3A1B")
	PathBg        = lipgloss.Color("#5B4634")
	PlayerFg      = lipgloss.Color("#722D4F")
	MarkFg        = lipgloss.Color("#90EE90")
	VisitedMarkFg = lipgloss.Color("#8EB173")
	VisitedPathBg = lipgloss.Color("#634E3A")
)
