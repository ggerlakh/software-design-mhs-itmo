package interpreter

import (
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
)

func TestInterpreter_Creation(t *testing.T) {
	env := map[string]string{
		"HOME": "/home/user",
		"PATH": "/usr/bin:/bin",
	}

	interpreter := &Interpreter{
		Env:       env,
		CmdParser: parser.Parser{},
	}

	if interpreter.Env["HOME"] != "/home/user" {
		t.Error("переменная окружения HOME не установлена правильно")
	}

	if interpreter.Env["PATH"] != "/usr/bin:/bin" {
		t.Error("переменная окружения PATH не установлена правильно")
	}

	if interpreter.CmdParser.BuiltinCommands != nil {
		t.Error("CmdParser.BuiltinCommands должен быть nil по умолчанию")
	}
}

func TestInterpreter_EmptyEnvironment(t *testing.T) {
	interpreter := &Interpreter{
		Env:       make(map[string]string),
		CmdParser: parser.Parser{},
	}

	if interpreter.Env == nil {
		t.Error("Env не должен быть nil")
	}

	if len(interpreter.Env) != 0 {
		t.Errorf("Env должен быть пустым, получено: %d элементов", len(interpreter.Env))
	}
}

func TestInterpreter_EnvironmentManipulation(t *testing.T) {
	interpreter := &Interpreter{
		Env:       make(map[string]string),
		CmdParser: parser.Parser{},
	}

	// Добавляем переменные окружения
	interpreter.Env["TEST_VAR"] = "test_value"
	interpreter.Env["ANOTHER_VAR"] = "another_value"

	if interpreter.Env["TEST_VAR"] != "test_value" {
		t.Error("переменная TEST_VAR не установлена правильно")
	}

	if interpreter.Env["ANOTHER_VAR"] != "another_value" {
		t.Error("переменная ANOTHER_VAR не установлена правильно")
	}

	// Изменяем переменную
	interpreter.Env["TEST_VAR"] = "modified_value"
	if interpreter.Env["TEST_VAR"] != "modified_value" {
		t.Error("переменная TEST_VAR не изменена правильно")
	}

	// Удаляем переменную
	delete(interpreter.Env, "ANOTHER_VAR")
	if _, exists := interpreter.Env["ANOTHER_VAR"]; exists {
		t.Error("переменная ANOTHER_VAR не удалена")
	}
}
