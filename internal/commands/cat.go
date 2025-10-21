package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// CatCommand реализует встроенную команду "cat".
// Она выводит содержимое файла (или stdin) на стандартный вывод.
type CatCommand struct{}

// Name возвращает имя команды.
func (c *CatCommand) Name() string {
	return "cat"
}

// Exec выполняет команду cat с переданными аргументами.
// Поддерживаются базовые опции:
//   - -n — нумеровать все строки
//   - -b — нумеровать только непустые строки (перекрывает -n)
//   - -s — убрать повторяющиеся пустые строки
//   - -E — показывать $ в конце строки
//   - -T — заменять табуляции на ^I
//
// Если не указано ни одного файла, читается stdin.
// Если указан "-" как имя файла, также читается stdin.
func (c *CatCommand) Exec(args []string) error {
	var (
		numberAll      bool
		numberNonEmpty bool
		squeezeBlank   bool
		showEnds       bool
		showTabs       bool
		files          []string
	)

	// Разбор опций
	for len(args) > 0 && strings.HasPrefix(args[0], "-") && len(args[0]) > 1 {
		switch args[0] {
		case "-n", "--number":
			numberAll = true
		case "-b", "--number-nonblank":
			numberNonEmpty = true
		case "-s", "--squeeze-blank":
			squeezeBlank = true
		case "-E", "--show-ends":
			showEnds = true
		case "-T", "--show-tabs":
			showTabs = true
		case "--help":
			fmt.Println(c.Help())
			return nil
		default:
			// Любая неизвестная опция считаем файлом
			files = append(files, args[0])
		}
		args = args[1:]
	}

	// Остаток аргументов считаем файлами
	if len(args) > 0 {
		files = append(files, args...)
	} else {
		files = []string{"-"}
	}

	lineNum := 1
	var prevBlank bool

	for _, fname := range files {
		var reader io.Reader

		if fname == "-" {
			reader = os.Stdin
		} else {
			//nolint:gosec // открываем файлы, как делает обычный cat, пользователь сам контролирует доступ
			file, err := os.Open(fname)
			if err != nil {
				return fmt.Errorf("cat: %v", err)
			}
			defer func(f *os.File) {
				if err := f.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "ошибка при закрытии файла %s: %v\n", f.Name(), err)
				}
			}(file)
			reader = file
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			isBlank := len(line) == 0

			if squeezeBlank && isBlank && prevBlank {
				continue
			}
			prevBlank = isBlank

			// Подстановка опций
			if showTabs {
				line = strings.ReplaceAll(line, "\t", "^I")
			}
			if showEnds {
				line += "$"
			}

			if numberNonEmpty {
				if !isBlank {
					fmt.Printf("%6d\t%s\n", lineNum, line)
					lineNum++
				} else {
					fmt.Println()
				}
			} else if numberAll {
				fmt.Printf("%6d\t%s\n", lineNum, line)
				lineNum++
			} else {
				fmt.Println(line)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("cat: %v", err)
		}
	}

	return nil
}

// Help возвращает справку по команде cat.
func (c *CatCommand) Help() string {
	return `NAME
    cat - объединяет содержимое файлов и выводит в стандартный вывод

SYNOPSIS
    cat [OPTION]... [FILE]...

DESCRIPTION
    Выводит содержимое указанных файлов в стандартный вывод.
    Если файл не указан или указан "-", читается стандартный ввод.

OPTIONS
    -b, --number-nonblank   нумеровать только непустые строки
    -n, --number             нумеровать все строки
    -s, --squeeze-blank      подавлять повторяющиеся пустые строки
    -E, --show-ends          отображать '$' в конце строки
    -T, --show-tabs          заменять TAB символом ^I
    --help                   показать эту справку

EXAMPLES
    cat file.txt
        → вывод содержимого file.txt

    cat -n file1.txt file2.txt
        → вывод содержимого file1.txt и file2.txt с нумерацией строк`
}
