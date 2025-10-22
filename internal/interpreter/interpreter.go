// Package interpreter реализует основной интерпретатор командной строки.
// Содержит структуру Interpreter для управления Read-Execute-Print Loop (REPL).
// Обрабатывает пользовательский ввод, парсит команды и выполняет их через пайплайны.
package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

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

// substitute выполняет подстановку переменных окружения в пользовательском вводе.
// Поддерживает подстановку переменных в формате $VAR или ${VAR}.
// Если переменная не найдена, оставляет строку как есть.
//
// Примеры:
//
//	echo $HOME → echo /home/user
//	echo ${PATH} → echo /usr/bin:/bin
//	echo $UNDEFINED → echo $UNDEFINED
func (i *Interpreter) substitute(userInput string) string {
	result := userInput

	// Подстановка переменных в формате ${VAR}
	result = regexp.MustCompile(`\$\{([^}]+)\}`).ReplaceAllStringFunc(result, func(match string) string {
		varName := match[2 : len(match)-1] // Убираем ${ и }
		if value, exists := i.Env[varName]; exists {
			return value
		}
		return match // Если переменная не найдена, оставляем как есть
	})

	// Подстановка переменных в формате $VAR
	result = regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`).ReplaceAllStringFunc(result, func(match string) string {
		varName := match[1:] // Убираем $
		if value, exists := i.Env[varName]; exists {
			return value
		}
		return match // Если переменная не найдена, оставляем как есть
	})

	return result
}
