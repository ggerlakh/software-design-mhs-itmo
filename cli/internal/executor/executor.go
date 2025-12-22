// Package executor отвечает за планирование и запуск команд,
// связывая их пайпами и управляя встроенными и внешними вызовами.
package executor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/ggerlakh/software-design-mhs-itmo/cli/internal/checkutils"
	"github.com/ggerlakh/software-design-mhs-itmo/cli/internal/commands"
)

// ExecutableCommand описывает команду, подготовленную к выполнению.
type ExecutableCommand struct {
	Name string
	Args []string
}

// Plan представляет последовательность команд, которые необходимо выполнить.
type Plan struct {
	Commands []ExecutableCommand
}

// Executor отвечает за выполнение команд согласно плану.
type Executor struct {
	BuiltinCommands []commands.BuiltinCommand
	Env             map[string]string
}

// NewExecutor создает новый Executor.
func NewExecutor(env map[string]string, builtins []commands.BuiltinCommand) *Executor {
	return &Executor{
		Env:             env,
		BuiltinCommands: builtins,
	}
}

// Execute запускает команды в соответствии с планом.
func (e *Executor) Execute(plan Plan) {
	if len(plan.Commands) == 0 {
		return
	}

	if len(plan.Commands) == 1 {
		ctx := e.newContext()
		ctx.Stdin = os.Stdin
		ctx.Stdout = os.Stdout
		ctx.Stderr = os.Stderr
		e.runCommand(plan.Commands[0], ctx)
		return
	}

	var pipes []*os.File

	defer func() {
		for _, pipe := range pipes {
			_ = pipe.Close()
		}
	}()

	contexts := make([]*commands.CommandContext, len(plan.Commands))
	for i := range contexts {
		contexts[i] = e.newContext()
	}

	for i := 0; i < len(plan.Commands)-1; i++ {
		reader, writer, err := os.Pipe()
		if err != nil {
			return
		}
		pipes = append(pipes, reader, writer)

		contexts[i].Stdout = writer
		contexts[i+1].Stdin = reader
	}

	contexts[0].Stdin = os.Stdin
	contexts[len(contexts)-1].Stdout = os.Stdout

	for i, cmd := range plan.Commands {
		if contexts[i].Stdout == nil {
			contexts[i].Stdout = os.Stdout
		}
		if contexts[i].Stderr == nil {
			contexts[i].Stderr = os.Stderr
		}
		e.runCommand(cmd, contexts[i])

		// Закрываем writer текущей команды.
		if writer, ok := contexts[i].Stdout.(*os.File); ok {
			_ = writer.Close()
		}
	}
}

func (e *Executor) newContext() *commands.CommandContext {
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}

	return &commands.CommandContext{
		Env: e.Env,
		Dir: currentDir,
	}
}

func (e *Executor) runCommand(cmd ExecutableCommand, ctx *commands.CommandContext) {
	switch {
	case checkutils.IsEnvAssignmentCommand(cmd.Name):
		parts := strings.SplitN(cmd.Name, "=", 2)
		e.Env[parts[0]] = parts[1]
	case checkutils.IsBuiltInCommand(cmd.Name, e.BuiltinCommands):
		for _, builtin := range e.BuiltinCommands {
			if builtin.Name() == cmd.Name {
				_ = builtin.Exec(cmd.Args, ctx)
			}
		}
	default:
		external := exec.Command(cmd.Name, cmd.Args...) //nolint:gosec
		external.Stdin = ctx.Stdin
		external.Stdout = ctx.Stdout
		external.Stderr = ctx.Stderr
		external.Dir = ctx.Dir

		for key, value := range ctx.Env {
			external.Env = append(external.Env, key+"="+value)
		}

		_ = external.Run()
	}
}
