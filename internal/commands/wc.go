package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// WcCommand реализует встроенную команду "wc".
// Она выводит количество строк, слов и байт в файле или stdin.
type WcCommand struct{}

// Name возвращает имя команды.
func (w *WcCommand) Name() string {
	return "wc"
}

// Exec выполняет команду wc с переданными аргументами.
// Если не указан файл, читается stdin.
//
// Примеры:
//
//	wc file.txt    → 3 10 55   (3 строки, 10 слов, 55 байт)
//	wc -l file.txt → 3         (только количество строк)
func (w *WcCommand) Exec(args []string) error {
	var files []string
	showLines, showWords, showBytes := true, true, true

	// Разбор опций
	for len(args) > 0 && strings.HasPrefix(args[0], "-") {
		switch args[0] {
		case "-l":
			showLines, showWords, showBytes = true, false, false
		case "-w":
			showLines, showWords, showBytes = false, true, false
		case "-c":
			showLines, showWords, showBytes = false, false, true
		}
		args = args[1:]
	}

	if len(args) == 0 {
		files = []string{"-"} // читаем stdin
	} else {
		files = args
	}

	for _, f := range files {
		var file *os.File
		var err error

		if f == "-" {
			file = os.Stdin
		} else {
			//nolint:gosec // открываем файлы, как делает обычный cat, пользователь сам контролирует доступ
			file, err = os.Open(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wc: не удалось открыть файл %s: %v\n", f, err)
				continue
			}
			defer func(f *os.File) {
				if err := f.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "ошибка при закрытии файла %s: %v\n", f.Name(), err)
				}
			}(file)
		}

		scanner := bufio.NewScanner(file)
		lines, words, bytesCount := 0, 0, 0

		for scanner.Scan() {
			line := scanner.Text()
			lines++
			words += len(strings.Fields(line))
			bytesCount += len(line) + 1 // +1 для \n
		}

		output := ""
		if showLines {
			output += fmt.Sprintf("%d ", lines)
		}
		if showWords {
			output += fmt.Sprintf("%d ", words)
		}
		if showBytes {
			output += fmt.Sprintf("%d ", bytesCount)
		}

		fmt.Println(strings.TrimSpace(output))
	}

	return nil
}

// Help возвращает справку по команде wc.
func (w *WcCommand) Help() string {
	return `NAME
    wc - подсчитывает количество строк, слов и байт в файле

SYNOPSIS
    wc [OPTION]... [FILE]...

DESCRIPTION
    Выводит количество строк, слов и байт в файле или stdin.

OPTIONS
    -l    показывать только количество строк
    -w    показывать только количество слов
    -c    показывать только количество байт

EXAMPLES
    wc file.txt
        → 3 10 55
    wc -l file.txt
        → 3`
}
