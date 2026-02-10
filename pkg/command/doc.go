// Package command provides the domain-specific language for player programs.
// It includes the lexer, parser, AST definitions, and interpreter.
package command

// This package is responsible for:
// - Lexical analysis of player code (tokenizer)
// - Parsing C-like syntax with Portuguese keywords
// - AST node definitions
// - Program execution/interpreter
// - Exposing to the application the actions
//
// Player-facing syntax examples:
//   repetir 3 { ← }
//   se parede { ↓ }
