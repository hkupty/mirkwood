package styles

import "github.com/charmbracelet/lipgloss"

// PathChar is the character used for filling path cells.
const (
	PathChar = ' '
)

// Walls style for maze wall cells (tree trunks/bark).
// TODO: Better fg color
// TODO: Figure out best symbol and decorations (i.e. borders instead of block)
var Walls = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#1B3A1B")).
	Background(lipgloss.Color("#2D2420"))

// Path style for empty path cells (forest floor).
// TODO: Better fg color
var Path = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#1B3A1B")).
	Background(lipgloss.Color("#4A3728"))

// Player style for the player character (bright highlight).
var Player = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFD700")).
	Background(lipgloss.Color("#4A3728"))

// Mark style for player-placed marks (distinct but subtle).
var Mark = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#90EE90")).
	Background(lipgloss.Color("#4A3728"))

// Visited style for cells the player has stepped on (subtle trail).
var Visited = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8B7355")).
	Background(lipgloss.Color("#4A3728"))

var VisitedMark = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#8EB173")).
	Background(lipgloss.Color("#4A3728"))
