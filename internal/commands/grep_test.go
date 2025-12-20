package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// testGrepExecWithOutput выполняет команду grep с переданными аргументами и возвращает stdout и stderr.
func testGrepExecWithOutput(cmd CommandExecutor, args []string, stdin io.Reader) (string, string) {
	var stdout, stderr bytes.Buffer
	if stdin == nil {
		stdin = strings.NewReader("")
	}
	ctx := &CommandContext{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}
	_ = cmd.Exec(args, ctx)
	return stdout.String(), stderr.String()
}

func TestGrepCommand_Name(t *testing.T) {
	cmd := &GrepCommand{}
	if cmd.Name() != "grep" {
		t.Errorf("ожидалось имя 'grep', получено %q", cmd.Name())
	}
}

func TestGrepCommand_Help(t *testing.T) {
	cmd := &GrepCommand{}
	help := cmd.Help()
	if !strings.Contains(help, "grep") {
		t.Error("справка должна содержать 'grep'")
	}
	if !strings.Contains(help, "-i") {
		t.Error("справка должна содержать описание флага -i")
	}
	if !strings.Contains(help, "-w") {
		t.Error("справка должна содержать описание флага -w")
	}
	if !strings.Contains(help, "-A") {
		t.Error("справка должна содержать описание флага -A")
	}
}

func TestGrepCommand_SimpleMatch(t *testing.T) {
	content := "hello world\nfoo bar\nhello again\n"
	tmp := t.TempDir() + "/test.txt"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"hello", tmp}, nil)

	expected := "hello world\nhello again\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_NoMatch(t *testing.T) {
	content := "hello world\nfoo bar\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"xyz", tmp}, nil)

	if out != "" {
		t.Errorf("ожидался пустой вывод, получено %q", out)
	}
}

func TestGrepCommand_CaseInsensitive(t *testing.T) {
	content := "Hello World\nHELLO AGAIN\nhello small\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-i", "hello", tmp}, nil)

	expected := "Hello World\nHELLO AGAIN\nhello small\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_WordMatch(t *testing.T) {
	content := "hello world\nhelloworld\nworld hello\nthehelloword\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-w", "hello", tmp}, nil)

	expected := "hello world\nworld hello\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_WordMatchWithUnicode(t *testing.T) {
	content := "привет мир\nприветмир\nмир привет\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-w", "привет", tmp}, nil)

	expected := "привет мир\nмир привет\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_AfterLines(t *testing.T) {
	content := "line1\nline2\nline3\nline4\nline5\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-A", "2", "line2", tmp}, nil)

	expected := "line2\nline3\nline4\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_AfterLinesZero(t *testing.T) {
	content := "line1\nline2\nline3\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-A", "0", "line2", tmp}, nil)

	expected := "line2\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_OverlappingAfterLines(t *testing.T) {
	// Тест на пересекающиеся области печати
	content := "match1\nafter1\nmatch2\nafter2\nend\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"-A", "2", "match", tmp}, nil)

	// match1 -> after1, match2
	// match2 -> after2, end
	// При пересечении строки печатаются только один раз
	expected := "match1\nafter1\nmatch2\nafter2\nend\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_RegexAnchorStart(t *testing.T) {
	content := "hello world\nworld hello\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"^hello", tmp}, nil)

	expected := "hello world\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_RegexAnchorEnd(t *testing.T) {
	content := "hello world\nworld hello\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"hello$", tmp}, nil)

	expected := "world hello\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_Stdin(t *testing.T) {
	cmd := &GrepCommand{}
	input := "line one\nline two\nline three\n"
	reader := strings.NewReader(input)
	out, _ := testGrepExecWithOutput(cmd, []string{"two"}, reader)

	expected := "line two\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_CombinedFlags(t *testing.T) {
	content := "Hello World\nHELLO there\nhelloWorld\ntest hello test\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	// -i -w: регистронезависимый поиск слова целиком
	out, _ := testGrepExecWithOutput(cmd, []string{"-i", "-w", "hello", tmp}, nil)

	expected := "Hello World\nHELLO there\ntest hello test\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_MissingPattern(t *testing.T) {
	cmd := &GrepCommand{}
	_, stderr := testGrepExecWithOutput(cmd, []string{}, nil)

	if !strings.Contains(stderr, "отсутствует паттерн") {
		t.Errorf("ожидалась ошибка об отсутствии паттерна, получено %q", stderr)
	}
}

func TestGrepCommand_InvalidRegex(t *testing.T) {
	cmd := &GrepCommand{}
	_, stderr := testGrepExecWithOutput(cmd, []string{"[invalid"}, nil)

	if !strings.Contains(stderr, "некорректное регулярное выражение") {
		t.Errorf("ожидалась ошибка о некорректном regex, получено %q", stderr)
	}
}

func TestGrepCommand_FileNotFound(t *testing.T) {
	cmd := &GrepCommand{}
	_, stderr := testGrepExecWithOutput(cmd, []string{"pattern", "/nonexistent/file.txt"}, nil)

	if !strings.Contains(stderr, "grep:") {
		t.Errorf("ожидалась ошибка grep:, получено %q", stderr)
	}
}

func TestGrepCommand_WordMatchAtBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		pattern  string
		expected string
	}{
		{
			name:     "word at start of line",
			content:  "word test\ntest word\n",
			pattern:  "word",
			expected: "word test\ntest word\n",
		},
		{
			name:     "word with underscore boundary",
			content:  "word_test\ntest_word\nword test\n",
			pattern:  "word",
			expected: "word test\n",
		},
		{
			name:     "word with number boundary",
			content:  "word123\n123word\nword test\n",
			pattern:  "word",
			expected: "word test\n",
		},
		{
			name:     "word with punctuation",
			content:  "word, test\ntest.word\nword\n",
			pattern:  "word",
			expected: "word, test\ntest.word\nword\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp := t.TempDir() + "/test.txt"
			_ = os.WriteFile(tmp, []byte(tt.content), 0o644)

			cmd := &GrepCommand{}
			out, _ := testGrepExecWithOutput(cmd, []string{"-w", tt.pattern, tmp}, nil)

			if out != tt.expected {
				t.Errorf("ожидалось %q, получено %q", tt.expected, out)
			}
		})
	}
}

func TestGrepCommand_AfterLinesAtEndOfFile(t *testing.T) {
	content := "line1\nline2\nmatch\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	// Запрашиваем 5 строк после, но в файле их меньше
	out, _ := testGrepExecWithOutput(cmd, []string{"-A", "5", "match", tmp}, nil)

	expected := "match\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_RegexSpecialChars(t *testing.T) {
	content := "file.txt\nfile_txt\nfiletxt\n"
	tmp := t.TempDir() + "/test.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &GrepCommand{}
	// Ищем точку (экранированную)
	out, _ := testGrepExecWithOutput(cmd, []string{"\\.txt", tmp}, nil)

	expected := "file.txt\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestGrepCommand_EmptyFile(t *testing.T) {
	tmp := t.TempDir() + "/empty.txt"
	_ = os.WriteFile(tmp, []byte(""), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"pattern", tmp}, nil)

	if out != "" {
		t.Errorf("ожидался пустой вывод для пустого файла, получено %q", out)
	}
}

func TestGrepCommand_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	file1 := dir + "/file1.txt"
	file2 := dir + "/file2.txt"
	_ = os.WriteFile(file1, []byte("hello from file1\n"), 0o644)
	_ = os.WriteFile(file2, []byte("hello from file2\n"), 0o644)

	cmd := &GrepCommand{}
	out, _ := testGrepExecWithOutput(cmd, []string{"hello", file1, file2}, nil)

	expected := "hello from file1\nhello from file2\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}
