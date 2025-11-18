package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/executor"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/preprocessor"
)

const exitCommand = "exit"

// Interpreter координирует работу препроцессинга, парсинга и выполнения команд.
type Interpreter struct {
	Preprocessor *preprocessor.Preprocessor
	Parser       *parser.Parser
	Executor     *executor.Executor
}

// Start запускает основной цикл интерпретатора (REPL).
func (i *Interpreter) Start() {
	fmt.Printf("Welcome to go-cli! To esacpe type %q.\n", exitCommand)
	scanner := bufio.NewScanner(os.Stdin)

Loop:
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break Loop
		}

		userInput := scanner.Text()

		preprocessed, err := i.Preprocessor.Process(userInput)
		if err != nil {
			fmt.Printf("preprocessing error: %s\n", err)
			continue
		}

		parsedPipeline, err := i.Parser.Parse(preprocessed)
		switch {
		case errors.Is(err, customErrors.ErrExit):
			break Loop
		case err != nil:
			fmt.Printf("%s\n", err)
			continue
		}

		executionPlan := toExecutionPlan(parsedPipeline)
		i.Executor.Execute(executionPlan)
	}
}

func toExecutionPlan(p parser.Pipeline) executor.Plan {
	plan := executor.Plan{
		Commands: make([]executor.ExecutableCommand, len(p.Commands)),
	}

	for idx, cmd := range p.Commands {
		plan.Commands[idx] = executor.ExecutableCommand{
			Name: cmd.Name,
			Args: append([]string{}, cmd.Args...),
		}
	}

	return plan
}
