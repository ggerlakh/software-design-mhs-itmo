// Package commands предоставляет интерфейсы и структуры для работы с командами.
// Содержит базовые интерфейсы CommandExecutor и BuiltinCommand, а также
// структуру CommandContext для передачи контекста выполнения команд.
package commands

import "io"

// CommandContext содержит контекст выполнения команды
type CommandContext struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    map[string]string
	Dir    string
}

// CommandExecutor определяет интерфейс для выполнения команд.
// Любая команда должна реализовывать метод Exec для выполнения
// с переданными аргументами и контекстом.
type CommandExecutor interface {
	Exec(args []string, ctx *CommandContext) error
}

// BuiltinCommand расширяет CommandExecutor для встроенных команд.
// Встроенные команды должны предоставлять имя и справку.
type BuiltinCommand interface {
	CommandExecutor
	Name() string
	Help() string
}
