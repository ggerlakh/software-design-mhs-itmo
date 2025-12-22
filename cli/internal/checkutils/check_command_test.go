package checkutils

import (
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/cli/internal/commands"
)

type stubCommand struct {
	name string
}

func (s *stubCommand) Name() string                                  { return s.name }
func (s *stubCommand) Help() string                                  { return "" }
func (s *stubCommand) Exec([]string, *commands.CommandContext) error { return nil }

func TestIsBuiltInCommand(t *testing.T) {
	builtins := []commands.BuiltinCommand{
		&stubCommand{name: "echo"},
		&stubCommand{name: "pwd"},
	}

	if !IsBuiltInCommand("echo", builtins) {
		t.Fatalf("ожидалось true для зарегистрированной команды")
	}

	if IsBuiltInCommand("unknown", builtins) {
		t.Fatalf("ожидалось false для неизвестной команды")
	}
}

func TestIsExternalCommand(t *testing.T) {
	if !IsExternalCommand("printf") {
		t.Fatalf("printf должен существовать в PATH")
	}

	if IsExternalCommand("command_that_does_not_exist_12345") {
		t.Fatalf("несуществующая команда не должна определяться как внешняя")
	}
}

func TestIsEnvAssignmentCommand(t *testing.T) {
	if !IsEnvAssignmentCommand("FOO=bar") {
		t.Fatalf("формат присваивания должен распознаваться")
	}

	if IsEnvAssignmentCommand("1INVALID=value") {
		t.Fatalf("некорректный формат не должен распознаваться")
	}
}
