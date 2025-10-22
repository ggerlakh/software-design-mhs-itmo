// Package interpreter реализует основной интерпретатор командной строки.
// Содержит структуру Interpreter для управления Read-Execute-Print Loop (REPL).
// Обрабатывает пользовательский ввод, парсит команды и выполняет их через пайплайны.
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

// Interpreter представляет основной интерпретатор командной строки.
// Реализует Read-Execute-Print Loop (REPL) для интерактивной работы.
type Interpreter struct {
	Env       map[string]string
	CmdParser parser.Parser
}

// Start запускает основной цикл интерпретатора (REPL).
// Читает команды из stdin, парсит их и выполняет.
// Завершается при команде "exit" или при EOF.
func (i *Interpreter) Start() {
	fmt.Printf("Welcome to go-cli! To esacpe type %q.\n", exitCommand)
	scanner := bufio.NewScanner(os.Stdin)
Loop:
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			// EOF или ошибка чтения - выходим из цикла
			break Loop
		}
		userInput := scanner.Text()
		pipeline, err := i.CmdParser.Parse(userInput, i.Env)

		switch {
		case errors.Is(err, customErrors.ErrExit):
			break Loop
		case err != nil:
			fmt.Printf("Evaluation error: %s\n", err)
		}

		pipeline.Run()
	}
}
