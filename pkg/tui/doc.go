// Package tui provides the terminal user interface for the Mirkwood game.
// It uses the bubbletea framework following the Elm Architecture pattern.
package tui

// This package is responsible for:
// - Rendering the maze display
// - Code input/editor interface
// - Execution visualization
// - Error display
// - Styling with lipgloss
//
// Import rules:
// - tui/ imports from core/ and maze/
// - core/ and maze/ must NOT import from tui/
