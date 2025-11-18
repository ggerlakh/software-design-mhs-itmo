package parser

import (
	"testing"

	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/preprocessor"
)

func newTestParser() *Parser {
	return NewParser([]string{"echo", "pwd", "cat", "wc"})
}

func TestParser_Parse_ExitCommand(t *testing.T) {
	parser := newTestParser()

	_, err := parser.Parse(preprocessor.Result{Original: "exit", Value: "exit"})

	if !customErrors.Is(err, customErrors.ErrExit) {
		t.Fatalf("ожидалась ошибка ErrExit, получено: %v", err)
	}
}

func TestParser_Parse_SingleCommand(t *testing.T) {
	parser := newTestParser()

	pipeline, err := parser.Parse(preprocessor.Result{Original: "echo hello", Value: "echo hello"})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 1 {
		t.Fatalf("ожидалась 1 команда, получено: %d", len(pipeline.Commands))
	}

	cmd := pipeline.Commands[0]
	if cmd.Name != "echo" || len(cmd.Args) != 1 || cmd.Args[0] != "hello" {
		t.Fatalf("неверно распознана команда: %#v", cmd)
	}
}

func TestParser_Parse_Pipeline(t *testing.T) {
	parser := newTestParser()

	pipeline, err := parser.Parse(preprocessor.Result{Original: "", Value: "echo hello | wc -w"})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 2 {
		t.Fatalf("ожидалось 2 команды, получено: %d", len(pipeline.Commands))
	}

	if pipeline.Commands[0].Name != "echo" {
		t.Fatalf("первая команда должна быть echo, получено: %s", pipeline.Commands[0].Name)
	}

	if pipeline.Commands[1].Name != "wc" {
		t.Fatalf("вторая команда должна быть wc, получено: %s", pipeline.Commands[1].Name)
	}
}

func TestParser_Parse_CommandNotFound(t *testing.T) {
	parser := newTestParser()

	_, err := parser.Parse(preprocessor.Result{Original: "", Value: "unknowncmd"})
	if err == nil {
		t.Fatalf("ожидалась ошибка для неизвестной команды")
	}
}

func TestParser_Parse_EmptyInput(t *testing.T) {
	parser := newTestParser()

	pipeline, err := parser.Parse(preprocessor.Result{Original: "", Value: "   "})
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 0 {
		t.Fatalf("для пустой строки ожидается 0 команд")
	}
}
