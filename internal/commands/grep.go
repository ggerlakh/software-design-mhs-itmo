// Package commands предоставляет реализации встроенных команд.

package commands

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"unicode"
)

// GrepCommand реализует встроенную команду "grep".
// Она ищет строки, соответствующие регулярному выражению в файле или stdin.
//
// Поддерживаемые флаги:
//   - -i — регистронезависимый поиск
//   - -w — поиск только слова целиком
//   - -A N — печатать N строк после совпадения
//
// Для разбора аргументов используется стандартная библиотека flag.
// Выбор обоснования:
//   - Не требует внешних зависимостей
//   - Достаточно мощная для поддержки необходимых флагов
//   - Хорошо документирована и широко используется в Go-проектах
//   - Поддерживает создание локального FlagSet для избежания конфликтов
type GrepCommand struct{}

// Name возвращает имя команды.
func (g *GrepCommand) Name() string {
	return "grep"
}

// grepFlags содержит распарсенные флаги для команды grep.
type grepFlags struct {
	ignoreCase bool // -i: регистронезависимый поиск
	wordMatch  bool // -w: поиск только слова целиком
	afterLines int  // -A: количество строк после совпадения
}

// parseGrepFlags разбирает аргументы командной строки для grep.
// Возвращает структуру с флагами, паттерн поиска и список файлов.
func parseGrepFlags(args []string) (*grepFlags, string, []string, error) {
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)

	// Создаём буфер для подавления вывода usage при ошибках парсинга
	fs.SetOutput(io.Discard)

	flags := &grepFlags{}
	fs.BoolVar(&flags.ignoreCase, "i", false, "регистронезависимый поиск")
	fs.BoolVar(&flags.wordMatch, "w", false, "поиск только слова целиком")
	fs.IntVar(&flags.afterLines, "A", 0, "количество строк после совпадения")

	if err := fs.Parse(args); err != nil {
		return nil, "", nil, fmt.Errorf("ошибка разбора флагов: %w", err)
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		return nil, "", nil, fmt.Errorf("grep: отсутствует паттерн для поиска")
	}

	pattern := remaining[0]
	files := remaining[1:]

	return flags, pattern, files, nil
}

// buildRegexp создаёт скомпилированное регулярное выражение на основе флагов.
// Если включён флаг -w, паттерн оборачивается в word boundaries.
// Если включён флаг -i, включается регистронезависимый режим.
func buildRegexp(pattern string, flags *grepFlags) (*regexp.Regexp, error) {
	// Примечание: для флага -w мы проверяем границы слов вручную в функции matchesWord,
	// так как \b в Go regexp не полностью поддерживает Unicode.
	// Паттерн не модифицируется здесь для -w.

	// Если включён регистронезависимый поиск, добавляем флаг (?i)
	if flags.ignoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: некорректное регулярное выражение: %w", err)
	}

	return re, nil
}

// isWordChar проверяет, является ли руна "word constituent character".
// Word constituent characters — это буквы, цифры и символ подчёркивания.
// Используется Unicode-классификация для поддержки не только ASCII.
func isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// matchesWord проверяет, соответствует ли найденное совпадение критерию "слово целиком".
// Совпадение считается словом, если оно ограничено не-word символами с обеих сторон.
func matchesWord(line string, match []int) bool {
	if len(match) < 2 {
		return false
	}

	start, end := match[0], match[1]

	// Проверяем символ перед совпадением
	if start > 0 {
		runes := []rune(line[:start])
		if len(runes) > 0 && isWordChar(runes[len(runes)-1]) {
			return false
		}
	}

	// Проверяем символ после совпадения
	if end < len(line) {
		runes := []rune(line[end:])
		if len(runes) > 0 && isWordChar(runes[0]) {
			return false
		}
	}

	return true
}

// lineMatches проверяет, соответствует ли строка регулярному выражению.
// Если включён флаг -w, проверяется также, что совпадение — слово целиком.
func lineMatches(line string, re *regexp.Regexp, wordMatch bool) bool {
	if !wordMatch {
		return re.MatchString(line)
	}

	// Для режима -w находим все совпадения и проверяем каждое
	matches := re.FindAllStringIndex(line, -1)
	for _, match := range matches {
		if matchesWord(line, match) {
			return true
		}
	}

	return false
}

// grepReader выполняет grep по содержимому reader и выводит результат в writer.
// Параметры:
//   - reader: источник данных для поиска
//   - writer: куда выводить результаты
//   - re: скомпилированное регулярное выражение
//   - flags: флаги команды grep
func grepReader(reader io.Reader, writer io.Writer, re *regexp.Regexp, flags *grepFlags) error {
	scanner := bufio.NewScanner(reader)

	// Сохраняем все строки для поддержки -A (after context)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("grep: ошибка чтения: %w", err)
	}

	// Множество для отслеживания уже напечатанных строк
	// (для обработки пересекающихся областей печати)
	printed := make(map[int]bool)

	for i, line := range lines {
		if lineMatches(line, re, flags.wordMatch) {
			// Печатаем совпавшую строку
			if !printed[i] {
				if _, err := fmt.Fprintln(writer, line); err != nil {
					return err
				}
				printed[i] = true
			}

			// Печатаем строки после совпадения (флаг -A)
			for j := 1; j <= flags.afterLines && i+j < len(lines); j++ {
				lineIdx := i + j
				if !printed[lineIdx] {
					if _, err := fmt.Fprintln(writer, lines[lineIdx]); err != nil {
						return err
					}
					printed[lineIdx] = true
				}
			}
		}
	}

	return nil
}

// Exec выполняет команду grep с переданными аргументами.
//
// Синтаксис:
//
//	grep [OPTIONS] PATTERN [FILE...]
//
// Опции:
//   - -i: регистронезависимый поиск
//   - -w: поиск только слова целиком
//   - -A N: печатать N строк после совпадения
//
// Если файлы не указаны, читается stdin.
//
// Примеры:
//
//	grep "hello" file.txt         → поиск "hello" в файле
//	grep -i "HELLO" file.txt      → регистронезависимый поиск
//	grep -w "word" file.txt       → поиск слова целиком
//	grep -A 2 "pattern" file.txt  → печать 2 строк после совпадения
func (g *GrepCommand) Exec(args []string, ctx *CommandContext) error {
	flags, pattern, files, err := parseGrepFlags(args)
	if err != nil {
		if _, writeErr := fmt.Fprintln(ctx.Stderr, err); writeErr != nil {
			return writeErr
		}
		return err
	}

	re, err := buildRegexp(pattern, flags)
	if err != nil {
		if _, writeErr := fmt.Fprintln(ctx.Stderr, err); writeErr != nil {
			return writeErr
		}
		return err
	}

	// Если файлы не указаны, читаем из stdin
	if len(files) == 0 {
		files = []string{"-"}
	}

	for _, fname := range files {
		var reader io.Reader

		if fname == "-" {
			reader = ctx.Stdin
		} else {
			//nolint:gosec // открываем файлы, как делает обычный grep
			file, err := os.Open(fname)
			if err != nil {
				if _, writeErr := fmt.Fprintf(ctx.Stderr, "grep: %v\n", err); writeErr != nil {
					return writeErr
				}
				continue
			}
			defer func(f *os.File) {
				if closeErr := f.Close(); closeErr != nil {
					if _, writeErr := fmt.Fprintf(ctx.Stderr,
						"grep: ошибка закрытия файла %s: %v\n", f.Name(), closeErr); writeErr != nil {
						// Игнорируем ошибку записи в stderr
						_ = writeErr
					}
				}
			}(file)
			reader = file
		}

		if err := grepReader(reader, ctx.Stdout, re, flags); err != nil {
			return err
		}
	}

	return nil
}

// Help возвращает справку по команде grep.
func (g *GrepCommand) Help() string {
	return `NAME
    grep — поиск строк, соответствующих регулярному выражению

SYNOPSIS
    grep [OPTIONS] PATTERN [FILE...]

DESCRIPTION
    Ищет строки, соответствующие регулярному выражению PATTERN,
    в указанных файлах или стандартном вводе.

OPTIONS
    -i          регистронезависимый поиск (case-insensitive)
    -w          поиск только слова целиком (word match)
    -A NUM      печатать NUM строк после каждого совпадения

    PATTERN     регулярное выражение для поиска
    FILE        файл(ы) для поиска; если не указаны, читается stdin

REGULAR EXPRESSIONS
    Поддерживаются регулярные выражения Go (RE2):
    .           любой символ
    *           0 или более повторений
    +           1 или более повторений
    ?           0 или 1 повторение
    ^           начало строки
    $           конец строки
    [abc]       любой из символов a, b, c
    [^abc]      любой символ, кроме a, b, c
    (...)       группировка
    |           альтернатива

EXAMPLES
    grep "error" log.txt
        → найти все строки с "error" в log.txt

    grep -i "ERROR" log.txt
        → регистронезависимый поиск "error"

    grep -w "test" file.txt
        → найти слово "test" целиком (не "testing")

    grep -A 3 "Exception" log.txt
        → найти "Exception" и 3 строки после

    grep "^#" config.txt
        → найти строки, начинающиеся с #

    grep "\.go$" files.txt
        → найти строки, заканчивающиеся на .go

NOTE ON -w FLAG
    Флаг -w ищет подстроки, ограниченные "non-word constituent character".
    Word constituent characters — это буквы (Unicode Letters), цифры
    (Unicode Digits) и символ подчёркивания (_).

NOTE ON OVERLAPPING REGIONS (-A flag)
    Если области печати (строки после совпадений) пересекаются,
    каждая строка печатается только один раз.

LIBRARY CHOICE
    Для разбора аргументов командной строки используется стандартная
    библиотека Go "flag". Выбор обоснован:
    
    Рассмотренные альтернативы:
    - pflag (github.com/spf13/pflag) — POSIX-совместимый, популярный
    - cobra (github.com/spf13/cobra) — полноценный CLI фреймворк
    - kong (github.com/alecthomas/kong) — современный, декларативный
    - go-flags (github.com/jessevdk/go-flags) — богатые возможности
    
    Выбрана стандартная библиотека "flag":
    ✓ Не требует внешних зависимостей
    ✓ Достаточно мощная для наших нужд (-i, -w, -A N)
    ✓ Поддерживает FlagSet для изоляции парсинга
    ✓ Хорошо документирована и стабильна
    ✓ Широко используется в Go-экосистеме`
}

// init проверяет, что GrepCommand реализует интерфейс BuiltinCommand.
// Это compile-time проверка.
var _ BuiltinCommand = (*GrepCommand)(nil)
