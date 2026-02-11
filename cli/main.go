// Package main is the application entry point for Mirkwood.
// It initializes the TUI and starts the game.
package main

import (
	"fmt"

	"github.com/hkupty/mirkwood/pkg/tui"
)

func main() {
	fmt.Println("Mirkwood - Educational Maze Game")
	fmt.Println("Run with: go run cli/main.go")

	tui.MainLoop()
}
