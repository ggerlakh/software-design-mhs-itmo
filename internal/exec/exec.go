package exec

import (
	"io"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
)

type Pipeline struct {
	Commands []ParsedCommand
}

func (p *Pipeline) Run() {
	panic("not implemented")
}

type ParsedCommand struct {
	Name            string
	Args            []string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
	Env             map[string]string
	CurrDir         string
	commandExecutor commands.CommandExecutor
}

func (p *ParsedCommand) Run() int {
	panic("not implemented")
}
