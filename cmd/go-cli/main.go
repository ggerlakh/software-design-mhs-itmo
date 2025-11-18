package main

import (
	"os"
	"strings"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/executor"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/interpreter"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/preprocessor"
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
	env := GetEnvMap()
	builtins := []commands.BuiltinCommand{
		&commands.CatCommand{},
		&commands.EchoCommand{},
		&commands.WcCommand{},
		&commands.PwdCommand{},
		&commands.ExitCommand{},
	}

	preproc := preprocessor.NewPreprocessor(&preprocessor.EnvSubstitutionStep{Env: env})
	cmdParser := parser.NewParser(builtinNames(builtins))
	exec := executor.NewExecutor(env, builtins)

	i := interpreter.Interpreter{
		Preprocessor: preproc,
		Parser:       cmdParser,
		Executor:     exec,
	}
	i.Start()
}

func builtinNames(cmds []commands.BuiltinCommand) []string {
	names := make([]string, len(cmds))
	for i, cmd := range cmds {
		names[i] = cmd.Name()
	}
	return names
}
