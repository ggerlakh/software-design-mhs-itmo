// Package errors содержит определения ошибок для интерпретатора командной строки.
// Предоставляет стандартные ошибки и функции для их проверки.
package errors

import (
	"errors"
	"fmt"
)

// ErrExit представляет ошибку завершения работы интерпретатора.
// Возвращается при выполнении команды "exit".
var ErrExit = errors.New("exit command")

// CommandNotFoundError представляет ошибку интерпретатора в случае если e.Command не была распознана
type CommandNotFoundError struct {
	Command string
}

func (e *CommandNotFoundError) Error() string {
	return fmt.Sprintf("go-cli: command not found: %s", e.Command)
}

// Is проверяет, является ли ошибка указанной ошибкой.
// Обертка над стандартной функцией errors.Is для удобства использования.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
