# Mirkwood — Project Context

## Agent Guidelines

It is of **utmost importance** that **NO CODE IS UPDATED UNLESS EXPLICITLY STATED**. It should never be **IMPLIED** that an action requests code updates. If ever in doubt, ask before taking action.

- **READING** code is fine
- **WRITING** or **MOVING** code, unless explicitly authorized, is prohibited

**Note:** By `code` we specifically mean `go.mod`, any file matching `*.go`, and/or any files within `cli/`, `mazes/`, or `pkg/` directories. Any other files are not subject to this rule.

## Overview

Mirkwood is a **TUI-based educational maze game** designed to teach children programming logic. Players write small algorithms to guide an agent through a maze, with progressive mechanics (loops, conditionals, mark, limits) introduced as levels advance.

**Language**: Go  
**TUI Stack**: bubbletea (Elm Architecture), lipgloss (styling), bubbles (components)

## Architecture

### Package Structure

```
├── cli/           # Application entry point (imports from pkg/tui/)
└── pkg/
    ├── core/      # Game logic, runtime state, validation (imports from maze/ and dsl/)
    ├── maze/      # Maze structure, generation, and win conditions
    ├── command/   # Player commands: lexer, parser, AST, and action transformations
    └── tui/       # Terminal interface (imports from core/)
```

### Dependency Rules

- `pkg/core/` imports from `pkg/maze/` and `pkg/command/`
- `pkg/tui/` imports from `pkg/core/` (adapter pattern)
- `cli/` imports from `pkg/tui/`
- Core logic (`pkg/maze/`, `pkg/core/`, `pkg/command/`) must NOT import TUI packages

### Maze as a Full Entity

A maze is a complete entity combining:
- **Static structure**: Grid layout (`MazeGrid`) and bitboard (`BitBoard`)
- **Win condition**: Defined in `maze.WinCondition` (exit, required marks, step limits)
- **Validation**: Checked by `core.Validator`

No separate `rules/` package - validation is split:
- `maze/` describes win conditions
- `core/` checks them at runtime

## Runtime Representation

- **Construction/Analysis**: Matrix/grid (`MazeGrid` - `[][]bool`)
- **Runtime/Execution**: 8×8 bitboards (`BitBoard` - `uint64`)

Bitboards (components in `core.State`):
- `Walls` — static, immutable (from `LevelBlueprint`)
- `Marks` — dynamic, player-placed
- `VisitedPath` — dynamic, tracks player movement

Runtime state (`core.State`):
- Agent position (single-bit bitboard)
- Marked bitboard
- VisitedPath bitboard
- Step counter
- Invariant view (read-only walls, finish point)

**Invariants:**
```
Marks & Walls == 0
VisitedPath & Walls == 0
Position & Walls == 0
```

## Execution Model

### Command Package Responsibilities

The `pkg/command/` package transforms player code into a **list of actions** (AST nodes). This is a pure transformation:
```
Parse(sourceCode) -> []Action | error
```

Actions represent executable operations like:
- `Move(North)`, `Move(South)`, etc.
- `ToggleMark()`
- `Turn(Left)`, `Turn(Right)`
- Control flow: `Repeat(n, actions)`, `If(condition, actions)`

### Core Responsibilities

The `pkg/core/` package **interprets** these actions by transforming the state sequentially:
```
Execute(actions, initialState) -> ExecutionResult
```

**Execution Rules:**

1. **Action Processing**: Each action transforms the state via pure functions:
   ```
   state.Move(dir) -> (newState, error)
   state.ToggleMark() -> newState
   ```

2. **Error Handling**: If any state transformation yields an error:
   - Halt execution immediately
   - Return the error
   - **Reset state** to the initial state (before execution began)

3. **Incomplete Solution**: If all actions are processed but `Position != FinishingPoint`:
   - Treat as error: "incomplete maze"
   - **Reset state** to the initial state

4. **Success Condition**: If at any point `Position == FinishingPoint`:
   - Halt execution immediately with success
   - Level is cleared
   - New level is loaded

5. **No Implicit Mutation**: State updates are always explicit and immutable

6. **Deterministic**: Same input always produces same output

## Player Language

**Internal**: AST (tree-based) defined in `pkg/command/`  
**Player-facing**: C-like block syntax with Portuguese keywords

Example:
```
repetir 3 {
  ←
}
```

**Player capabilities:**
- Sense local environment (wall ahead)
- Mark/sense current cell
- Move and turn

**Engine-only (not exposed):**
- Visited cell knowledge
- Adjacency queries
- Pathfinding
- Win condition checking

## Testing

- Unit tests for core logic
- Fuzz tests where applicable (parsing, bitboards, invariants)
- Tests for invariants
- Tests for execution semantics
- Deterministic tests via seeded randomness
