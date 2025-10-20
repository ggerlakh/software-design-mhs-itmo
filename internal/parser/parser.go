package parser

import (
	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/exec"
)

type Parser struct {
	BuiltinCommands []commands.BuiltinCommand
}

func (p *Parser) Parse(substitutedInput string, globalEnv map[string]string) (exec.Pipeline, error) {
	panic("not implemented")
}
