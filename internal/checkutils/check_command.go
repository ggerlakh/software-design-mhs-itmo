// Package checkutils содержит вспомогательные функции для проверки команд:
// builtin, внешних и присваиваний переменных окружения.
package checkutils

import (
	stdExec "os/exec"
	"regexp"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
)

func IsBuiltInCommand(commandName string, builtinCommands []commands.BuiltinCommand) bool {
	for _, cmd := range builtinCommands {
		if cmd.Name() == commandName {
			return true
		}
	}
	return false
}

func IsExternalCommand(commandName string) bool {
	_, err := stdExec.LookPath(commandName)
	return err == nil
}

func IsEnvAssignmentCommand(command string) bool {
	// Регулярное выражение для проверки формата VAR=value
	envPattern := `^[a-zA-Z_][a-zA-Z0-9_]*=.*$`
	matched, _ := regexp.MatchString(envPattern, command)
	return matched
}
