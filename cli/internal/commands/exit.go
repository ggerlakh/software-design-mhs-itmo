package commands

import (
	"github.com/ggerlakh/software-design-mhs-itmo/cli/internal/errors"
)

type ExitCommand struct{}

func (e *ExitCommand) Name() string {
	return "exit"
}

func (e *ExitCommand) Exec(args []string, ctx *CommandContext) error {
	return errors.ErrExit
}

func (e *ExitCommand) Help() string {
	return "exit - terminate the shell"
}
