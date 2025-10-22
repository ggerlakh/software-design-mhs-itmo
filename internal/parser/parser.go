// Package parser отвечает за разбор пользовательского ввода на команды и пайпы.
// Преобразует строки команд в структуры ParsedCommand и Pipeline для последующего выполнения.
// Поддерживает одиночные команды, пайпы и команду exit для завершения работы.
package parser

import (
	"os"
	"strings"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/exec"
)

// Parser отвечает за разбор пользовательского ввода на команды и пайпы.
// Содержит список встроенных команд для проверки при разборе.
type Parser struct {
	BuiltinCommands []commands.BuiltinCommand
}

// Parse разбирает входную строку на команды и создает Pipeline.
// Поддерживает:
//   - Одиночные команды: "echo hello"
//   - Пайпы: "echo hello | wc"
//   - Команду exit для завершения работы
//
// Возвращает Pipeline с разобранными командами или ошибку.
func (p *Parser) Parse(substitutedInput string, globalEnv map[string]string) (exec.Pipeline, error) {
	// Проверяем на команду exit
	if strings.TrimSpace(substitutedInput) == "exit" {
		return exec.Pipeline{}, customErrors.ErrExit
	}

	// Разбиваем входную строку по символу пайпа
	commandStrings := strings.Split(substitutedInput, "|")

	var comms []exec.ParsedCommand

	for _, cmdStr := range commandStrings {
		cmdStr = strings.TrimSpace(cmdStr)
		if cmdStr == "" {
			continue
		}

		// Разбираем команду на имя и аргументы
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		commandName := parts[0]
		args := parts[1:]

		// Получаем текущую директорию
		currentDir, err := os.Getwd()
		if err != nil {
			// Если не удалось получить текущую директорию, используем "."
			currentDir = "."
		}

		// Создаем CommandContext
		ctx := &commands.CommandContext{
			Env: globalEnv,
			Dir: currentDir,
		}

		// Создаем ParsedCommand
		parsedCmd := exec.ParsedCommand{
			Name:    commandName,
			Args:    args,
			Context: ctx,
		}

		comms = append(comms, parsedCmd)
	}

	return exec.Pipeline{Commands: comms, BuiltinCommands: p.BuiltinCommands}, nil
}
