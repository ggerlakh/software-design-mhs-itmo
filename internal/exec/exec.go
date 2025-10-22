// Package exec отвечает за выполнение команд и пайпов.
// Содержит структуры Pipeline и ParsedCommand для управления выполнением команд.
// Поддерживает как встроенные команды, так и внешние программы с настройкой потоков ввода/вывода.
package exec

import (
	"os"
	"os/exec"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
)

// Pipeline представляет последовательность команд, соединенных пайпами.
// Содержит список команд для выполнения и встроенные команды для проверки.
type Pipeline struct {
	Commands        []ParsedCommand
	BuiltinCommands []commands.BuiltinCommand
}

// Run выполняет все команды в пайплайне.
// Для одиночных команд настраивает stdin/stdout/stderr напрямую.
// Для множественных команд создает пайпы между ними и выполняет последовательно.
func (p *Pipeline) Run() {
	if len(p.Commands) == 0 {
		return
	}

	// Если команда одна, выполняем её напрямую
	if len(p.Commands) == 1 {
		// Настраиваем stdin/stdout для одиночной команды
		p.Commands[0].Context.Stdin = os.Stdin
		p.Commands[0].Context.Stdout = os.Stdout
		p.Commands[0].Context.Stderr = os.Stderr
		p.Commands[0].Run(p.BuiltinCommands)
		return
	}

	// Для множественных команд создаем пайпы
	var pipes []*os.File
	defer func() {
		// Закрываем все пайпы
		for _, pipe := range pipes {
			if err := pipe.Close(); err != nil {
				// Игнорируем ошибки закрытия пайпов
				_ = err
			}
		}
	}()

	// Создаем пайпы между командами
	for i := 0; i < len(p.Commands)-1; i++ {
		reader, writer, err := os.Pipe()
		if err != nil {
			return
		}
		pipes = append(pipes, reader, writer)

		// Настраиваем stdin для следующей команды
		p.Commands[i+1].Context.Stdin = reader
		// Настраиваем stdout для текущей команды
		p.Commands[i].Context.Stdout = writer
	}

	// Настраиваем stdin для первой команды и stdout для последней
	p.Commands[0].Context.Stdin = os.Stdin
	p.Commands[len(p.Commands)-1].Context.Stdout = os.Stdout

	// Настраиваем stderr для всех команд
	for i := range p.Commands {
		p.Commands[i].Context.Stderr = os.Stderr
	}

	// Выполняем команды последовательно для правильной работы пайпов
	for i := range p.Commands {
		p.Commands[i].Run(p.BuiltinCommands)

		// Закрываем writer после выполнения команды
		if i < len(p.Commands)-1 {
			// Закрываем writer для текущей команды
			if writer, ok := p.Commands[i].Context.Stdout.(*os.File); ok {
				if err := writer.Close(); err != nil {
					// Игнорируем ошибки закрытия writer
					_ = err
				}
			}
		}
	}
}

// ParsedCommand представляет разобранную команду с именем, аргументами и контекстом.
type ParsedCommand struct {
	Name    string
	Args    []string
	Context *commands.CommandContext
}

// Run выполняет команду, сначала проверяя встроенные команды,
// затем пытаясь выполнить как внешнюю программу.
// Возвращает код возврата: 0 для успеха, 1 для ошибки.
func (p *ParsedCommand) Run(builtinCommands []commands.BuiltinCommand) int {
	// Ищем встроенную команду
	for _, builtin := range builtinCommands {
		if builtin.Name() == p.Name {
			err := builtin.Exec(p.Args, p.Context)
			if err != nil {
				return 1
			}
			return 0
		}
	}

	// Если не найдена встроенная команда, пытаемся выполнить внешнюю
	//nolint:gosec // Пользователь сам контролирует выполнение команд
	cmd := exec.Command(p.Name, p.Args...)
	cmd.Stdin = p.Context.Stdin
	cmd.Stdout = p.Context.Stdout
	cmd.Stderr = p.Context.Stderr
	cmd.Dir = p.Context.Dir

	// Устанавливаем переменные окружения
	for key, value := range p.Context.Env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	err := cmd.Run()
	if err != nil {
		return 1
	}
	return 0
}
