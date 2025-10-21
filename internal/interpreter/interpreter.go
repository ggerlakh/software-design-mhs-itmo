package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	customErrors "github.com/ggerlakh/software-design-mhs-itmo/internal/errors"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
)

const exitCommand = "exit"

type Interpreter struct {
	Env       map[string]string
	CmdParser parser.Parser
}

// Точка запуска интерпретатора
func (i *Interpreter) Start() {
	fmt.Printf("Welcome to go-cli! To esacpe type %q.\n", exitCommand)
	scanner := bufio.NewScanner(os.Stdin)
Loop:
	for {
		fmt.Print("> ")
		scanner.Scan()
		userInput := scanner.Text()
		substitutedInput := i.substitute(userInput)
		pipeline, err := i.CmdParser.Parse(substitutedInput, i.Env)

		switch {
		case errors.Is(err, customErrors.ErrExit):
			break Loop
		case err != nil:
			fmt.Printf("Evaluation error: %s\n", err)
		}

		pipeline.Run()
	}
}

//nolint:unused
func (i *Interpreter) substitute(userInput string) string {
	// TODO
	return userInput
}
