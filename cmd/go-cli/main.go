package main

import (
	"os"
	"strings"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/interpreter"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
)

func GetEnvMap() map[string]string {
	envMap := make(map[string]string)

	for _, env := range os.Environ() {
		// Разделяем строку на ключ и значение по первому '='
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}

	return envMap
}

func main() {
	i := interpreter.Interpreter{
		Env: GetEnvMap(),
		CmdParser: parser.Parser{
			BuiltinCommands: []commands.BuiltinCommand{
				&commands.CatCommand{},
				&commands.EchoCommand{},
				&commands.WcCommand{},
				&commands.PwdCommand{},
				&commands.ExitCommand{},
			},
		},
	}
	i.Start()
}
