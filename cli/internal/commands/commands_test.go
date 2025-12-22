package commands

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// TestCommandContext создает тестовый контекст команды
func TestCommandContext_Creation(t *testing.T) {
	env := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin:/bin",
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env,
		Dir:    "/tmp",
	}

	if ctx.Stdin != os.Stdin {
		t.Error("Stdin не установлен правильно")
	}

	if ctx.Stdout != os.Stdout {
		t.Error("Stdout не установлен правильно")
	}

	if ctx.Stderr != os.Stderr {
		t.Error("Stderr не установлен правильно")
	}

	if ctx.Dir != "/tmp" {
		t.Errorf("Dir не установлен правильно, ожидалось: /tmp, получено: %s", ctx.Dir)
	}

	if ctx.Env["HOME"] != "/home/user" {
		t.Error("переменная окружения HOME не установлена правильно")
	}

	if ctx.Env["PATH"] != "/usr/bin:/bin" {
		t.Error("переменная окружения PATH не установлена правильно")
	}
}

func TestCommandContext_EmptyEnvironment(t *testing.T) {
	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	if ctx.Env == nil {
		t.Error("Env не должен быть nil")
	}

	if len(ctx.Env) != 0 {
		t.Errorf("Env должен быть пустым, получено: %d элементов", len(ctx.Env))
	}
}

func TestCommandContext_WithBuffers(t *testing.T) {
	var stdinBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	ctx := &CommandContext{
		Stdin:  &stdinBuf,
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	// Тестируем запись в stdout
	testOutput := "test output\n"
	n, err := ctx.Stdout.Write([]byte(testOutput))
	if err != nil {
		t.Errorf("ошибка записи в stdout: %v", err)
	}

	if n != len(testOutput) {
		t.Errorf("ожидалось записать %d байт, записано: %d", len(testOutput), n)
	}

	if stdoutBuf.String() != testOutput {
		t.Errorf("ожидался вывод %q, получено: %q", testOutput, stdoutBuf.String())
	}

	// Тестируем запись в stderr
	testError := "test error\n"
	_, err = ctx.Stderr.Write([]byte(testError))
	if err != nil {
		t.Errorf("ошибка записи в stderr: %v", err)
	}

	if stderrBuf.String() != testError {
		t.Errorf("ожидался вывод ошибки %q, получено: %q", testError, stderrBuf.String())
	}
}

func TestCommandContext_EnvironmentManipulation(t *testing.T) {
	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	// Добавляем переменные окружения
	ctx.Env["TEST_VAR"] = "test_value"
	ctx.Env["ANOTHER_VAR"] = "another_value"

	if ctx.Env["TEST_VAR"] != "test_value" {
		t.Error("переменная TEST_VAR не установлена правильно")
	}

	if ctx.Env["ANOTHER_VAR"] != "another_value" {
		t.Error("переменная ANOTHER_VAR не установлена правильно")
	}

	// Изменяем переменную
	ctx.Env["TEST_VAR"] = "modified_value"
	if ctx.Env["TEST_VAR"] != "modified_value" {
		t.Error("переменная TEST_VAR не изменена правильно")
	}

	// Удаляем переменную
	delete(ctx.Env, "ANOTHER_VAR")
	if _, exists := ctx.Env["ANOTHER_VAR"]; exists {
		t.Error("переменная ANOTHER_VAR не удалена")
	}
}

// TestCommandExecutor интерфейс для тестирования
type TestCommandExecutor struct {
	name     string
	execFunc func(args []string, ctx *CommandContext) error
}

func (t *TestCommandExecutor) Exec(args []string, ctx *CommandContext) error {
	if t.execFunc != nil {
		return t.execFunc(args, ctx)
	}
	return nil
}

func TestCommandExecutor_Interface(t *testing.T) {
	var output bytes.Buffer

	executor := &TestCommandExecutor{
		name: "test",
		execFunc: func(args []string, ctx *CommandContext) error {
			ctx.Stdout.Write([]byte("test command executed\n"))
			return nil
		},
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	err := executor.Exec([]string{"arg1", "arg2"}, ctx)
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	expected := "test command executed\n"
	if output.String() != expected {
		t.Errorf("ожидался вывод %q, получено: %q", expected, output.String())
	}
}

func TestCommandExecutor_WithError(t *testing.T) {
	var output bytes.Buffer

	executor := &TestCommandExecutor{
		name: "test",
		execFunc: func(args []string, ctx *CommandContext) error {
			ctx.Stderr.Write([]byte("error occurred\n"))
			return os.ErrPermission
		},
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: &output,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	err := executor.Exec([]string{}, ctx)
	if err == nil {
		t.Error("ожидалась ошибка")
	}

	if err != os.ErrPermission {
		t.Errorf("ожидалась ошибка os.ErrPermission, получено: %v", err)
	}

	if output.String() != "error occurred\n" {
		t.Errorf("ожидался вывод ошибки, получено: %q", output.String())
	}
}

// TestBuiltinCommand для тестирования интерфейса BuiltinCommand
type TestBuiltinCommand struct {
	name     string
	helpText string
	execFunc func(args []string, ctx *CommandContext) error
}

func (t *TestBuiltinCommand) Name() string {
	return t.name
}

func (t *TestBuiltinCommand) Exec(args []string, ctx *CommandContext) error {
	if t.execFunc != nil {
		return t.execFunc(args, ctx)
	}
	return nil
}

func (t *TestBuiltinCommand) Help() string {
	return t.helpText
}

func TestBuiltinCommand_Interface(t *testing.T) {
	var output bytes.Buffer

	builtin := &TestBuiltinCommand{
		name:     "testbuiltin",
		helpText: "This is a test builtin command",
		execFunc: func(args []string, ctx *CommandContext) error {
			ctx.Stdout.Write([]byte("builtin executed\n"))
			return nil
		},
	}

	// Тестируем методы интерфейса
	if builtin.Name() != "testbuiltin" {
		t.Errorf("ожидалось имя 'testbuiltin', получено: %s", builtin.Name())
	}

	if builtin.Help() != "This is a test builtin command" {
		t.Errorf("ожидалась справка 'This is a test builtin command', получено: %s", builtin.Help())
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	err := builtin.Exec([]string{"arg1"}, ctx)
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	expected := "builtin executed\n"
	if output.String() != expected {
		t.Errorf("ожидался вывод %q, получено: %q", expected, output.String())
	}
}

func TestBuiltinCommand_EmptyArgs(t *testing.T) {
	var output bytes.Buffer

	builtin := &TestBuiltinCommand{
		name:     "test",
		helpText: "test help",
		execFunc: func(args []string, ctx *CommandContext) error {
			if len(args) != 0 {
				t.Errorf("ожидалось 0 аргументов, получено: %d", len(args))
			}
			ctx.Stdout.Write([]byte("no args\n"))
			return nil
		},
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	err := builtin.Exec([]string{}, ctx)
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	expected := "no args\n"
	if output.String() != expected {
		t.Errorf("ожидался вывод %q, получено: %q", expected, output.String())
	}
}

func TestBuiltinCommand_MultipleArgs(t *testing.T) {
	var output bytes.Buffer

	builtin := &TestBuiltinCommand{
		name:     "test",
		helpText: "test help",
		execFunc: func(args []string, ctx *CommandContext) error {
			expectedArgs := []string{"arg1", "arg2", "arg3"}
			if len(args) != len(expectedArgs) {
				t.Errorf("ожидалось %d аргументов, получено: %d", len(expectedArgs), len(args))
			}

			for i, expected := range expectedArgs {
				if args[i] != expected {
					t.Errorf("аргумент %d: ожидалось %s, получено: %s", i, expected, args[i])
				}
			}

			ctx.Stdout.Write([]byte("args processed\n"))
			return nil
		},
	}

	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &output,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	err := builtin.Exec([]string{"arg1", "arg2", "arg3"}, ctx)
	if err != nil {
		t.Errorf("неожиданная ошибка: %v", err)
	}

	expected := "args processed\n"
	if output.String() != expected {
		t.Errorf("ожидался вывод %q, получено: %q", expected, output.String())
	}
}

func TestCommandContext_IOOperations(t *testing.T) {
	var stdinBuf bytes.Buffer
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	// Записываем данные в stdin
	testInput := "test input data\n"
	stdinBuf.Write([]byte(testInput))

	ctx := &CommandContext{
		Stdin:  &stdinBuf,
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
		Env:    make(map[string]string),
		Dir:    ".",
	}

	// Читаем из stdin
	buf := make([]byte, 1024)
	n, err := ctx.Stdin.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("ошибка чтения из stdin: %v", err)
	}

	readData := string(buf[:n])
	if readData != testInput {
		t.Errorf("ожидалось прочитать %q, прочитано: %q", testInput, readData)
	}

	// Записываем в stdout
	testOutput := "test output data\n"
	_, err = ctx.Stdout.Write([]byte(testOutput))
	if err != nil {
		t.Errorf("ошибка записи в stdout: %v", err)
	}

	if stdoutBuf.String() != testOutput {
		t.Errorf("ожидался вывод %q, получено: %q", testOutput, stdoutBuf.String())
	}

	// Записываем в stderr
	testError := "test error data\n"
	_, err = ctx.Stderr.Write([]byte(testError))
	if err != nil {
		t.Errorf("ошибка записи в stderr: %v", err)
	}

	if stderrBuf.String() != testError {
		t.Errorf("ожидался вывод ошибки %q, получено: %q", testError, stderrBuf.String())
	}
}
