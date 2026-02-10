# Mirkwood — TODO

This document tracks the next steps for the Mirkwood educational maze game project.

## Current State

### Existing Code
- **pkg/maze/model.go**: Core maze types (BitBoard, MazeGrid, Direction, LevelBlueprint, WinCondition)
- **pkg/maze/conversion.go**: Grid/BitBoard conversion utilities
- **pkg/core/state.go**: Runtime State with bitboard operations
- **pkg/core/movement.go**: Move and ToggleMark operations with errors
- **pkg/core/blueprint.go**: State initialization from LevelBlueprint
- **pkg/core/validation.go**: Win condition checking (Validator)

### What's Working
- 8×8 maze representation with bitboards
- Immutable State transitions (Move, ToggleMark)
- Wall collision detection
- Direction-based movement (N/S/E/W)
- Win conditions (reach exit, required marks, step limits)
- State validation (invariant enforcement)
- New package structure: pkg/core/, pkg/maze/, pkg/command/, pkg/tui/, cli/

## Immediate Tasks

### Core Logic
- [ ] **Maze Generation**: Implement maze generation algorithm (recursive backtracker or Prim's)
  - Generate valid 8×8 mazes with guaranteed solution paths
  - Ensure walls wrap edges as frame
- [ ] **Level Persistence**: Load/save LevelBlueprint from files
- [ ] **Sensors**: Implement player-facing sensor functions
  - `wallAhead()` - check wall in current direction
  - `isMarked()` - check if current cell is marked

### Command (pkg/command/)
- [ ] **Action Types**: Define action types (Move, Turn, ToggleMark, etc.)
- [ ] **AST Definition**: Define AST nodes for control flow (Repeat, If)
- [ ] **Lexer**: Tokenize player code (Portuguese keywords, arrows)
- [ ] **Parser**: Transform code into **list of actions** (AST)
  - `repetir N { ... }` for loops
  - `se condicao { ... }` for conditionals
  - Arrow symbols (← ↑ → ↓) for movement
  - Returns `[]Action` or parse error
- [ ] **Unfolding**: Flatten/unfold control flow into linear action list where needed

### Execution (pkg/core/)
- [ ] **Action Interpreter**: Process list of actions against State
  - Sequential action processing
  - Halt on first error (return error + reset state)
  - Halt on reaching FinishingPoint (success, load next level)
  - Halt on completion without reaching FinishingPoint (incomplete maze error + reset state)
  - Pure function: `Execute(actions, state) -> ExecutionResult`
- [ ] **Restructure Entry Points**: Make `Step` the only public entry point for TUI
  - `Process` should be internal/testing-only
  - TUI calls `Step` individually with timing control (200-800ms delays)
  - Remove timing/display concerns from core execution logic

### TUI (pkg/tui/)
- [ ] **Bubble Tea Setup**: Initialize TUI framework
- [ ] **Maze Rendering**: Display maze with walls, agent position, marks
- [ ] **Code Editor**: Input area for player programs
- [ ] **Execution Visualization**: Step-through animation of player code
- [ ] **Error Display**: User-friendly error messages for syntax/runtime errors

### Testing
- [ ] **Fuzz Tests**: Add fuzzing where applicable
  - Bitboard operations (shifts, masks, collisions)
  - Maze generation (valid mazes have solutions)
  - State transitions (Move in all directions)
  - Parser (valid programs parse correctly)
  - Move `Process` function to test utilities for action sequence fuzzing
- [ ] **Invariant Tests**: Property-based tests for bitboard invariants
- [ ] **Deterministic Tests**: Seeded randomness for maze generation tests

## Design Decisions Pending

### State Management
- State is immutable (good!)
- Consider: Should invariants be exported or kept internal? (Currently internal with view)
- Consider: How to represent direction the player is facing? (Currently not tracked)

### Player Language Evolution
- Level 1: Basic movement (arrows only)
- Level 2: Add loops (`repetir`)
- Level 3: Add conditionals (`se`)
- Level 4: Add `mark` capability
- Later: Introduce step limits, no-revisit constraints

## Read Further

**Architecture:**
- *Game Engine Architecture* by Jason Gregory — ECS patterns
- *Data-Oriented Design* by Richard Fabian — Cache efficiency and bitboards

**Go Fuzzing:**
- https://go.dev/doc/security/fuzz/
- https://go.dev/doc/tutorial/fuzz

**Maze Generation:**
- Recursive Backtracker algorithm
- Prim's algorithm for mazes
- Wilson's algorithm (uniform spanning tree)

**TUI Development:**
- https://github.com/charmbracelet/bubbletea
- https://github.com/charmbracelet/lipgloss
- Elm Architecture pattern

## Red Flags to Watch For

- [ ] Core logic importing TUI packages
- [ ] Player language exposing engine internals
- [ ] Bitboard logic leaking into maze generation
- [ ] Implicit state mutation
- [ ] Missing invariant enforcement
- [ ] String errors instead of typed errors
- [ ] Missing fuzz tests for parsing/bitboards
- [ ] `utils/` packages shared across domains
