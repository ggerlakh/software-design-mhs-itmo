// Package parser отвечает за разбор пользовательского ввода на команды и пайпы.
// Преобразует результат препроцессинга в независимую модель ParsedPipeline.
// Поддерживает одиночные команды, пайпы и команду exit для завершения работы.
package parser

import (
	"strings"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/checkutils"
	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/preprocessor"
)

const exitCommand = "exit"

// ParsedCommand описывает команду, полученную после парсинга.
type ParsedCommand struct {
	Name string
	Args []string
}

// Pipeline представляет результат парсинга командной строки.
type Pipeline struct {
	Commands []ParsedCommand
}

// Parser отвечает за валидацию и разбор пользовательского ввода.
// Хранит список известных встроенных команд для проверки.
type Parser struct {
	builtinNames map[string]struct{}
}

// NewParser создает парсер с перечнем имен встроенных команд.
func NewParser(builtinNames []string) *Parser {
	nameSet := make(map[string]struct{}, len(builtinNames))
	for _, name := range builtinNames {
		nameSet[name] = struct{}{}
	}

	return &Parser{builtinNames: nameSet}
}

// Parse превращает результат препроцессинга в Pipeline.
func (p *Parser) Parse(input preprocessor.Result) (Pipeline, error) {
	if strings.TrimSpace(input.Value) == "" {
		return Pipeline{}, nil
	}

	if strings.TrimSpace(input.Value) == exitCommand {
		return Pipeline{}, customErrors.ErrExit
	}

	segments := strings.Split(input.Value, "|")
	var commands []ParsedCommand

	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		parts := strings.Fields(segment)
		if len(parts) == 0 {
			continue
		}

		name := parts[0]
		args := parts[1:]

		if !p.isKnownCommand(name) &&
			!checkutils.IsExternalCommand(name) &&
			!checkutils.IsEnvAssignmentCommand(name) {
			return Pipeline{}, &customErrors.CommandNotFoundError{Command: name}
		}

		commands = append(commands, ParsedCommand{
			Name: name,
			Args: args,
		})
	}

	return Pipeline{Commands: commands}, nil
}

func (p *Parser) isKnownCommand(name string) bool {
	_, ok := p.builtinNames[name]
	return ok
}
