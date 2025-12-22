package commands

import (
	"fmt"
	"os"
)

// PwdCommand реализует встроенную команду "pwd".
// Она выводит текущий рабочий каталог.
type PwdCommand struct{}

// Name возвращает имя команды.
func (p *PwdCommand) Name() string {
	return "pwd"
}

// Exec выполняет команду pwd с переданными аргументами.
// Команда игнорирует все аргументы и просто выводит текущий путь.
//
// Примеры:
//
//	pwd → /home/user/project
func (p *PwdCommand) Exec(args []string, ctx *CommandContext) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(ctx.Stdout, dir); err != nil {
		return err
	}
	return nil
}

// Help возвращает справку по команде pwd.
func (p *PwdCommand) Help() string {
	return `NAME
    pwd - выводит текущий рабочий каталог

SYNOPSIS
    pwd

DESCRIPTION
    Печатает путь к текущему рабочему каталогу.

EXAMPLES
    pwd
        → /home/user/project`
}
