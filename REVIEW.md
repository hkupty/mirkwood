# Mirkwood — Code Review Guide for Local Agent

## Purpose of This Document

This document defines **how a local review agent should evaluate the codebase** for this project.  
Its goal is not stylistic enforcement, but **architectural correctness, pedagogical integrity, and conceptual soundness**.

The reviewer should assume that:
- This is a learning-oriented project
- Clarity, invariants, and boundaries matter more than cleverness
- Performance is explored deliberately, not prematurely
- The project is intended to evolve incrementally

The reviewer should **challenge design decisions** when they violate stated principles, even if the code “works”.

---

## Project Summary (Context)

This project is a **TUI-based educational maze game** designed to teach children programming logic.

Players write small algorithms to guide an agent through a maze.
As levels progress, new mechanics and constraints are introduced (loops, conditionals, mark, limits).

The implementation language is **Go**.

The TUI is built using:
- **bubbletea** — The Elm Architecture framework for Go
- **lip gloss** — Style definitions for terminal applications
- **bubbles** — Pre-built TUI components
- **harmonica** *(optional)* — Smooth animations and transitions

The system is explicitly divided into:
- A **pure logic core** (engine)
- A **TUI frontend**

The codebase should reflect this separation clearly.

---

## High-Level Design Principles

The reviewer should verify that the code adheres to these principles:

### 1. Pedagogy first
- Teaching logic is the primary goal
- Engine conveniences must not leak into the player language
- No hidden “solver shortcuts” exposed to players

### 2. Explicit boundaries
- Core logic must not depend on IO
- Player-visible capabilities must be strictly limited
- Engine-only knowledge must remain internal

### 3. Representations are intentional
- High-level representations for construction/analysis
- Low-level representations for runtime execution
- Conversion between the two is explicit and one-way

### 4. Invariants over validation
- Illegal states should be unrepresentable where possible
- Bitwise masking and construction-time guarantees are preferred over runtime checks

---

## Architecture Expectations

### Layering

The directory structure is intentionally flexible — **the code is the source of truth**. Rather than enforcing a specific folder hierarchy, the reviewer should verify that the architectural boundaries are respected:

- The **core logic** (maze, state, interpreter, rules) must remain pure and free of IO dependencies
- The **TUI layer** (rendering, input, views) must act as a thin adapter over the core
- Core packages must not import TUI-specific packages

Key rules:
- Core logic must not perform IO or depend on terminal libraries
- TUI components should delegate all game logic to the core
- State transitions and game rules are owned by the core, not the UI

---

## Maze Representation Rules

### Construction / Analysis Phase

Used for:
- Maze generation
- Pathfinding
- Validation
- Difficulty analysis

Expected representation:
- Matrix / grid (`[][]Cell`, `[][]bool`, or equivalent)

Reviewer checks:
- Maze generation algorithms are readable and conventional
- Pathfinding is implemented clearly (BFS/A*, not bitwise hacks)
- This representation is not reused for runtime execution

---

### Runtime / Evaluation Phase

Used for:
- Executing player programs
- Enforcing constraints
- Validating outcomes

Expected representation:
- 8×8 **bitboards** (`uint64`)

Required bitboards:
- `Walls` (static)
- `Marked` (dynamic)
- `VisitedPath` (visited cells, dynamic)

---

### Bitboard Invariants

The reviewer must verify these **hard invariants** are enforced:

```

Marked & Walls == 0
VisitedPath & Walls == 0

```

The walls wrap the edges of the map as a frame, preventing position over/underflow to adjacent rows.

Enforcement expectations:
- Wall collisions are prevented at move time
- Mark and path updates are masked or constrained
- Illegal states should not be reachable

The reviewer should flag:
- Any runtime code that allows overlapping wall/mark/path bits
- Any logic that relies on post-hoc validation instead of prevention

---

## Runtime State Model

Expected minimal runtime state:

- Agent position (single cell index)
- Marked bitboard
- VisitedPath bitboard

Reviewer checks:
- Maze (walls) is immutable at runtime
- Mark and path belong to **state**, not maze
- State transitions are explicit and deterministic

---

## Execution Semantics

### Core rule

All player commands must be executed via a single, explicit mechanism, conceptually:

```

apply(op, state, maze) -> newState | error

```

Reviewer checks:
- No command mutates state implicitly
- Failure modes are explicit (wall hit, step limit, invalid command)
- Execution is deterministic and replayable

---

## Player vs Engine Capabilities (Critical Boundary)

### Engine-only capabilities (allowed)

The engine may:
- Know which cells were visited
- Inspect adjacency and graph structure
- Perform pathfinding
- Analyze difficulty
- Generate hints

These capabilities **must not** be exposed to the player language.

---

### Player-visible capabilities (strictly limited)

The player language may only:
- Sense local environment (e.g. wall ahead)
- Mark the current cell
- Sense mark on the current cell
- Move and turn

Reviewer must flag:
- Any exposure of `isVisited`
- Any adjacency or graph queries
- Any API that allows implicit memory or global inspection

Memory and strategy must be **explicitly constructed** using mark.

---

## Command Language & Syntax

### Internal Representation

- Commands compile to an **AST**
- Tree-based, Lisp-like semantics internally
- Interpreter operates only on AST

Reviewer checks:
- No logic depends on surface syntax
- AST is the single source of truth for execution

---

### Player-Facing Syntax

- **C-like block syntax**
- Visual nesting with braces
- Symbols for actions (← ↑ → ↓, mark)
- Portuguese keywords for control flow (`repetir`, `se`)

Example:
```

repetir 3 {
←
}

```

Reviewer checks:
- Syntax is designed for readability, not parser convenience
- Error messages are understandable for non-experts
- Grammar complexity increases gradually

---

## VisitedPath Bitboard Usage

The `VisitedPath` (visited) bitboard is expected to be used for:

- Validation rules (no revisits, coverage)
- Replay and visualization
- Difficulty analysis

Reviewer must ensure:
- VisitedPath is **not** used to detect collisions (that happens at move time)
- VisitedPath knowledge is **not** exposed to player logic

---

## Performance Expectations

Performance is a **secondary concern**, used as a learning exercise.

Reviewer should verify:
- No premature optimization in construction algorithms
- Bitboards are used only in runtime execution
- Profiling or benchmarking code is isolated and optional

---

## Error Handling & Feedback

Reviewer checks:
- Errors are categorized (syntax, runtime, rule violation)
- Failure modes are explicit and explainable
- Engine errors do not crash the CLI

---

## Testing Expectations

Reviewer should expect:
- Unit tests for core logic
- Tests for invariants
- Tests for execution semantics
- Deterministic tests via seeded randomness

Lack of tests is acceptable early, but architecture must be testable.

---

## Red Flags (Must Be Called Out)

The reviewer should explicitly flag:

- Player access to engine knowledge
- Mixing IO with core logic
- Bitboard logic leaking into generation/pathfinding
- Overly clever representations harming clarity
- Implicit state mutation
- Hidden solvers or shortcuts

---

## Reviewer Mindset

The reviewer should act as:
- An architecture reviewer
- A pedagogy guardian
- A correctness auditor

Not as:
- A style enforcer
- A micro-optimizer
- A framework evangelist

---

## Final Guiding Question

For any piece of code, the reviewer should ask:

> “Does this make the system easier to reason about, harder to misuse, and better at teaching logic?”

If the answer is “no”, it should be challenged.
```
