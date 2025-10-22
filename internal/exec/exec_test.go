package exec

import (
	"bytes"
	"os"
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
)

func TestPipeline_Run_EmptyPipeline(t *testing.T) {
	pipeline := &Pipeline{
		Commands:        []ParsedCommand{},
		BuiltinCommands: []commands.BuiltinCommand{},
	}

	// Не должно паниковать
	pipeline.Run()
}

func TestParsedCommand_Run_BuiltinCommand(t *testing.T) {
	var output bytes.Buffer

	mockCmd := &commands.EchoCommand{}

	ctx := &commands.CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	parsedCmd := ParsedCommand{
		Name:    "echo",
		Args:    []string{"hello", "world"},
		Context: ctx,
	}

	builtinCommands := []commands.BuiltinCommand{mockCmd}

	exitCode := parsedCmd.Run(builtinCommands)

	if exitCode != 0 {
		t.Errorf("ожидался код выхода 0, получено: %d", exitCode)
	}

	expected := "hello world\n"
	if output.String() != expected {
		t.Errorf("ожидался вывод %q, получено %q", expected, output.String())
	}
}

func TestParsedCommand_Run_UnknownCommand(t *testing.T) {
	var output bytes.Buffer

	ctx := &commands.CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: &output,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	parsedCmd := ParsedCommand{
		Name:    "nonexistent",
		Args:    []string{},
		Context: ctx,
	}

	builtinCommands := []commands.BuiltinCommand{}

	exitCode := parsedCmd.Run(builtinCommands)

	// Должен вернуть код ошибки, так как команда не найдена
	if exitCode != 1 {
		t.Errorf("ожидался код выхода 1 для неизвестной команды, получено: %d", exitCode)
	}
}

func TestParsedCommand_Run_EnvironmentVariables(t *testing.T) {
	var output bytes.Buffer

	env := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin:/bin",
	}

	ctx := &commands.CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: &output,
		Env:    env,
		Dir:    ".",
	}

	parsedCmd := ParsedCommand{
		Name:    "nonexistent",
		Args:    []string{},
		Context: ctx,
	}

	builtinCommands := []commands.BuiltinCommand{}

	exitCode := parsedCmd.Run(builtinCommands)

	// Проверяем, что переменные окружения переданы
	if ctx.Env["HOME"] != "/home/user" {
		t.Errorf("переменная окружения HOME не передана правильно")
	}
	if ctx.Env["PATH"] != "/usr/bin:/bin" {
		t.Errorf("переменная окружения PATH не передана правильно")
	}

	if exitCode != 1 {
		t.Errorf("ожидался код выхода 1 для неизвестной команды, получено: %d", exitCode)
	}
}

func TestParsedCommand_Run_PwdCommand(t *testing.T) {
	var output bytes.Buffer

	mockCmd := &commands.PwdCommand{}

	ctx := &commands.CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	parsedCmd := ParsedCommand{
		Name:    "pwd",
		Args:    []string{},
		Context: ctx,
	}

	builtinCommands := []commands.BuiltinCommand{mockCmd}

	exitCode := parsedCmd.Run(builtinCommands)

	if exitCode != 0 {
		t.Errorf("ожидался код выхода 0, получено: %d", exitCode)
	}

	// Проверяем, что вывод содержит текущую директорию
	if output.String() == "" {
		t.Error("ожидался вывод текущей директории")
	}
}
