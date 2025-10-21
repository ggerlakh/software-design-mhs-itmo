package commands

import (
	"fmt"
	"strings"
)

// EchoCommand реализует встроенную команду "echo".
// Она выводит переданные аргументы в стандартный вывод,
// разделяя их пробелами. По умолчанию добавляет перевод строки в конце.
type EchoCommand struct{}

// Name возвращает имя команды.
func (e *EchoCommand) Name() string {
	return "echo"
}

// Exec выполняет команду echo с переданными аргументами.
// Поддерживаются базовые опции:
//   - -n — не добавлять перевод строки в конце вывода.
//
// Примеры:
//
//	echo hello world   → hello world\n
//	echo -n test       → test
func (e *EchoCommand) Exec(args []string) error {
	newline := true

	// Проверяем опцию -n
	if len(args) > 0 && args[0] == "-n" {
		newline = false
		args = args[1:]
	}

	// Собираем строку из аргументов, разделяя пробелами
	output := strings.Join(args, " ")

	// Печатаем результат с учётом опции -n
	if newline {
		fmt.Println(output)
	} else {
		fmt.Print(output)
	}

	return nil
}

// Help возвращает справку по команде echo.
func (e *EchoCommand) Help() string {
	return `NAME
    echo - выводит переданные аргументы в стандартный вывод

SYNOPSIS
    echo [OPTION]... [STRING]...

DESCRIPTION
    Печатает переданные строки, разделённые пробелами.
    По умолчанию добавляет перевод строки в конце.

OPTIONS
    -n  — не добавлять перевод строки в конце вывода

EXAMPLES
    echo Hello world
        → Hello world

    echo -n Hello
        → Hello`
}
