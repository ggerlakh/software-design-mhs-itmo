package commands

import (
	"strings"
	"testing"

	customErrors "github.com/ggerlakh/software-design-mhs-itmo/cli/internal/errors"
)

func TestBuiltinCommandMetadata(t *testing.T) {
	tests := []struct {
		name     string
		cmd      BuiltinCommand
		expected string
	}{
		{"cat", &CatCommand{}, "cat"},
		{"echo", &EchoCommand{}, "echo"},
		{"wc", &WcCommand{}, "wc"},
		{"pwd", &PwdCommand{}, "pwd"},
		{"exit", &ExitCommand{}, "exit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.Name() != tt.expected {
				t.Fatalf("ожидалось Name=%s, получено %s", tt.expected, tt.cmd.Name())
			}

			help := tt.cmd.Help()
			if help == "" {
				t.Fatalf("Help не должен быть пустым")
			}
			if !strings.Contains(strings.ToLower(help), tt.name) {
				t.Fatalf("Help должен содержать название команды, help=%q", help)
			}
		})
	}
}

func TestExitCommandExecReturnsErrExit(t *testing.T) {
	cmd := &ExitCommand{}
	err := cmd.Exec(nil, &CommandContext{})
	if !customErrors.Is(err, customErrors.ErrExit) {
		t.Fatalf("ожидался ErrExit, получено %v", err)
	}
}
