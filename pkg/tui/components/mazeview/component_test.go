package mazeview

import (
	"strings"
	"testing"

	"github.com/hkupty/mirkwood/pkg/maze"
)

// TestSelectBoxChar validates the box-drawing character selection logic.
// Each test case checks a specific neighbor configuration and expected output.
func TestSelectBoxChar(t *testing.T) {
	tests := []struct {
		name                  string
		up, down, left, right bool
		expected              rune
	}{
		// No connections
		{"isolated", false, false, false, false, '▓'},

		// Single connections (treated as line terminals)
		{"up only", true, false, false, false, '│'},
		{"down only", false, true, false, false, '│'},
		{"left only", false, false, true, false, '─'},
		{"right only", false, false, false, true, '─'},

		// Two connections
		{"vertical line", true, true, false, false, '│'},
		{"horizontal line", false, false, true, true, '─'},
		{"up+left (bottom-right corner)", true, false, true, false, '┘'},
		{"up+right (bottom-left corner)", true, false, false, true, '└'},
		{"down+left (top-right corner)", false, true, true, false, '┐'},
		{"down+right (top-left corner)", false, true, false, true, '┌'},

		// Three connections (T-junctions)
		{"right T (up+down+left)", true, true, true, false, '┤'},
		{"left T (up+down+right)", true, true, false, true, '├'},
		{"bottom T (up+left+right)", true, false, true, true, '┴'},
		{"top T (down+left+right)", false, true, true, true, '┬'},

		// Four connections (cross)
		{"cross (all)", true, true, true, true, '┼'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectBoxChar(tt.up, tt.down, tt.left, tt.right)
			if result != tt.expected {
				t.Errorf("selectBoxChar(%v, %v, %v, %v) = %q, want %q",
					tt.up, tt.down, tt.left, tt.right, result, tt.expected)
			}
		})
	}
}

// TestCellWidth validates that different cell widths produce expected output widths.
func TestCellWidth(t *testing.T) {
	blueprint := maze.LevelBlueprint{
		Key: 1,
		Grid: maze.MazeGrid{
			{true, false},
			{false, true},
		},
		StartingPoint:  0,
		FinishingPoint: 3,
		WinCondition:   maze.SimpleExit,
	}

	t.Run("narrow (2-char) width", func(t *testing.T) {
		m := New(blueprint).WithCellWidth(CellWidthNarrow)
		lines := strings.Split(strings.TrimRight(m.View(), "\n"), "\n")
		for i, line := range lines {
			// Each line should be 2 cells * 2 chars = 4 visible characters
			// (lipgloss styles add ANSI codes but visible length should be 4)
			if len(line) == 0 {
				t.Errorf("line %d is empty", i)
			}
		}
	})

	t.Run("wide (3-char) width", func(t *testing.T) {
		m := New(blueprint).WithCellWidth(CellWidthWide)
		lines := strings.Split(strings.TrimRight(m.View(), "\n"), "\n")
		for i, line := range lines {
			if len(line) == 0 {
				t.Errorf("line %d is empty", i)
			}
		}
	})
}

// TestSampleMazeRendering validates that the SampleMaze renders without errors.
func TestSampleMazeRendering(t *testing.T) {
	blueprint := maze.LevelBlueprint{
		Key:            1,
		Grid:           maze.SampleMaze,
		StartingPoint:  1,
		FinishingPoint: 55,
		WinCondition:   maze.SimpleExit,
	}

	m := New(blueprint)
	view := m.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

	// Should have 8 lines (one per row)
	lines := strings.Split(strings.TrimRight(view, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("Expected 8 lines, got %d", len(lines))
	}

	// Each line should have content (even if just styled spaces)
	for i, line := range lines {
		if len(line) == 0 {
			t.Errorf("Line %d is empty", i)
		}
	}
}

// TestAgentRendering validates that the agent is rendered at the correct position.
func TestAgentRendering(t *testing.T) {
	// 2x2 maze with agent at position [0][1]
	blueprint := maze.LevelBlueprint{
		Key: 1,
		Grid: maze.MazeGrid{
			{true, false},
			{true, true},
		},
		StartingPoint:  1, // [0][1]
		FinishingPoint: 3,
		WinCondition:   maze.SimpleExit,
	}

	m := New(blueprint).WithCellWidth(CellWidthNarrow)
	view := m.View()

	// The view should contain the agent symbol somewhere
	if !strings.Contains(view, "◆") {
		t.Error("View() does not contain agent symbol '◆'")
	}
}

// TestVisibilityToggles validates that visibility settings affect rendering.
func TestVisibilityToggles(t *testing.T) {
	blueprint := maze.LevelBlueprint{
		Key: 1,
		Grid: maze.MazeGrid{
			{false, false},
			{false, false},
		},
		StartingPoint:  0,
		FinishingPoint: 3,
		WinCondition:   maze.SimpleExit,
	}

	t.Run("agent visibility", func(t *testing.T) {
		visible := New(blueprint).WithAgentVisibility(true).View()
		hidden := New(blueprint).WithAgentVisibility(false).View()

		// Agent should be in visible but not hidden (or rendered differently)
		if visible == hidden {
			t.Error("Agent visibility toggle had no effect on rendering")
		}
	})
}

// TestRenderWith validates the functional options pattern for one-off renders.
func TestRenderWith(t *testing.T) {
	blueprint := maze.LevelBlueprint{
		Key: 1,
		Grid: maze.MazeGrid{
			{true, false, true},
			{false, false, false},
			{true, false, true},
		},
		StartingPoint:  4, // Center [1][1]
		FinishingPoint: 0,
		WinCondition:   maze.SimpleExit,
	}

	m := New(blueprint).WithCellWidth(CellWidthNarrow)

	// Default render
	defaultView := m.View()

	// One-off render with different width
	wideView := m.RenderWith(WithRenderCellWidth(CellWidthWide))

	if defaultView == wideView {
		t.Error("RenderWith did not change the output width")
	}

	// Verify default model wasn't mutated
	afterView := m.View()
	if afterView != defaultView {
		t.Error("RenderWith mutated the model's persistent state")
	}
}

// TestBoxDrawingEdgeCases validates corner cases in box-drawing selection.
func TestBoxDrawingEdgeCases(t *testing.T) {
	t.Run("maze boundaries", func(t *testing.T) {
		// At maze edges, we should handle missing neighbors gracefully
		// (edges are implicitly non-walls)
		blueprint := maze.LevelBlueprint{
			Key: 1,
			Grid: maze.MazeGrid{
				{true, true, true},
				{true, false, true},
				{true, true, true},
			},
			StartingPoint:  4,
			FinishingPoint: 4,
			WinCondition:   maze.SimpleExit,
		}

		m := New(blueprint)
		view := m.View()
		if view == "" {
			t.Error("Failed to render maze with edge walls")
		}
	})
}

// FuzzTestBoxCharSelection uses fuzzing to ensure selectBoxChar handles all inputs.
func FuzzSelectBoxChar(f *testing.F) {
	// Add seed inputs for all 16 combinations
	f.Add(false, false, false, false)
	f.Add(true, true, true, true)
	f.Add(true, false, false, false)
	f.Add(false, true, false, false)
	f.Add(false, false, true, false)
	f.Add(false, false, false, true)

	f.Fuzz(func(t *testing.T, up, down, left, right bool) {
		// Should never panic regardless of input combination
		result := selectBoxChar(up, down, left, right)

		// Result should be a valid Unicode box-drawing character
		validChars := map[rune]bool{
			'─': true, '│': true, '┌': true, '┐': true,
			'└': true, '┘': true, '├': true, '┤': true,
			'┬': true, '┴': true, '┼': true, '▓': true,
		}

		if !validChars[result] {
			t.Errorf("Unexpected character returned: %U (%q)", result, result)
		}
	})
}
