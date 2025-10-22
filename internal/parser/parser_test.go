package parser

import (
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
)

func TestParser_Parse_ExitCommand(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("exit", env)

	if !customErrors.Is(err, customErrors.ErrExit) {
		t.Errorf("ожидалась ошибка ErrExit, получено: %v", err)
	}

	if len(pipeline.Commands) != 0 {
		t.Errorf("для команды exit не должно быть команд в пайпе, получено: %d", len(pipeline.Commands))
	}
}

func TestParser_Parse_ExitCommandWithSpaces(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("  exit  ", env)

	if !customErrors.Is(err, customErrors.ErrExit) {
		t.Errorf("ожидалась ошибка ErrExit для команды с пробелами, получено: %v", err)
	}

	if len(pipeline.Commands) != 0 {
		t.Errorf("для команды exit не должно быть команд в пайпе, получено: %d", len(pipeline.Commands))
	}
}

func TestParser_Parse_SingleCommand(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("echo hello world", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 1 {
		t.Errorf("ожидалась 1 команда, получено: %d", len(pipeline.Commands))
	}

	cmd := pipeline.Commands[0]
	if cmd.Name != "echo" {
		t.Errorf("ожидалось имя команды 'echo', получено: %s", cmd.Name)
	}

	expectedArgs := []string{"hello", "world"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("ожидалось %d аргументов, получено: %d", len(expectedArgs), len(cmd.Args))
	}

	for i, expected := range expectedArgs {
		if cmd.Args[i] != expected {
			t.Errorf("аргумент %d: ожидалось %s, получено: %s", i, expected, cmd.Args[i])
		}
	}
}

func TestParser_Parse_CommandWithSpaces(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("  echo   hello   world  ", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 1 {
		t.Errorf("ожидалась 1 команда, получено: %d", len(pipeline.Commands))
	}

	cmd := pipeline.Commands[0]
	if cmd.Name != "echo" {
		t.Errorf("ожидалось имя команды 'echo', получено: %s", cmd.Name)
	}

	expectedArgs := []string{"hello", "world"}
	if len(cmd.Args) != len(expectedArgs) {
		t.Errorf("ожидалось %d аргументов, получено: %d", len(expectedArgs), len(cmd.Args))
	}
}

func TestParser_Parse_Pipeline(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("echo hello | wc -w", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 2 {
		t.Errorf("ожидалось 2 команды в пайпе, получено: %d", len(pipeline.Commands))
	}

	// Первая команда
	cmd1 := pipeline.Commands[0]
	if cmd1.Name != "echo" {
		t.Errorf("первая команда: ожидалось имя 'echo', получено: %s", cmd1.Name)
	}
	expectedArgs1 := []string{"hello"}
	if len(cmd1.Args) != len(expectedArgs1) {
		t.Errorf("первая команда: ожидалось %d аргументов, получено: %d", len(expectedArgs1), len(cmd1.Args))
	}

	// Вторая команда
	cmd2 := pipeline.Commands[1]
	if cmd2.Name != "wc" {
		t.Errorf("вторая команда: ожидалось имя 'wc', получено: %s", cmd2.Name)
	}
	expectedArgs2 := []string{"-w"}
	if len(cmd2.Args) != len(expectedArgs2) {
		t.Errorf("вторая команда: ожидалось %d аргументов, получено: %d", len(expectedArgs2), len(cmd2.Args))
	}
}

func TestParser_Parse_PipelineWithSpaces(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("  echo hello  |  wc -w  ", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 2 {
		t.Errorf("ожидалось 2 команды в пайпе, получено: %d", len(pipeline.Commands))
	}

	cmd1 := pipeline.Commands[0]
	if cmd1.Name != "echo" {
		t.Errorf("первая команда: ожидалось имя 'echo', получено: %s", cmd1.Name)
	}

	cmd2 := pipeline.Commands[1]
	if cmd2.Name != "wc" {
		t.Errorf("вторая команда: ожидалось имя 'wc', получено: %s", cmd2.Name)
	}
}

func TestParser_Parse_EmptyString(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 0 {
		t.Errorf("для пустой строки не должно быть команд, получено: %d", len(pipeline.Commands))
	}
}

func TestParser_Parse_OnlySpaces(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("   ", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 0 {
		t.Errorf("для строки только с пробелами не должно быть команд, получено: %d", len(pipeline.Commands))
	}
}

func TestParser_Parse_EmptyCommands(t *testing.T) {
	parser := &Parser{}
	env := make(map[string]string)

	pipeline, err := parser.Parse("echo hello || wc", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	// Должно быть 3 команды: echo hello, пустая команда, wc
	if len(pipeline.Commands) != 2 {
		t.Errorf("ожидалось 2 команды (пустые команды должны игнорироваться), получено: %d", len(pipeline.Commands))
	}

	cmd1 := pipeline.Commands[0]
	if cmd1.Name != "echo" {
		t.Errorf("первая команда: ожидалось имя 'echo', получено: %s", cmd1.Name)
	}

	cmd2 := pipeline.Commands[1]
	if cmd2.Name != "wc" {
		t.Errorf("вторая команда: ожидалось имя 'wc', получено: %s", cmd2.Name)
	}
}

func TestParser_Parse_EnvironmentVariables(t *testing.T) {
	parser := &Parser{}
	env := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin:/bin",
	}

	pipeline, err := parser.Parse("echo test", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.Commands) != 1 {
		t.Errorf("ожидалась 1 команда, получено: %d", len(pipeline.Commands))
	}

	cmd := pipeline.Commands[0]
	if cmd.Context.Env["HOME"] != "/home/user" {
		t.Errorf("переменная окружения HOME не установлена правильно")
	}
	if cmd.Context.Env["PATH"] != "/usr/bin:/bin" {
		t.Errorf("переменная окружения PATH не установлена правильно")
	}
}

func TestParser_Parse_BuiltinCommands(t *testing.T) {
	builtinCommands := []commands.BuiltinCommand{
		&commands.EchoCommand{},
		&commands.PwdCommand{},
	}
	parser := &Parser{BuiltinCommands: builtinCommands}
	env := make(map[string]string)

	pipeline, err := parser.Parse("echo test", env)

	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	if len(pipeline.BuiltinCommands) != 2 {
		t.Errorf("ожидалось 2 встроенные команды, получено: %d", len(pipeline.BuiltinCommands))
	}
}
