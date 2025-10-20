package interpreter

import (
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
)

type Interpreter struct {
	Env       map[string]string
	CmdParser parser.Parser
}

func (i *Interpreter) Start() {
	panic("not implemented")
}

func (i *Interpreter) substitue(userInput string) string {
	panic("not implemented")
}
