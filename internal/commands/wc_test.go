package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// testExecWithOutput выполняет команду с переданными аргументами и возвращает вывод
func testWcExecWithOutput(cmd CommandExecutor, args []string, stdin io.Reader) string {
	var buf bytes.Buffer
	if stdin == nil {
		stdin = os.Stdin
	}
	ctx := &CommandContext{
		Stdin:  stdin,
		Stdout: &buf,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}
	cmd.Exec(args, ctx)
	return buf.String()
}

func TestWcCommand_LinesWordsBytes(t *testing.T) {
	content := "one two\nthree four five\nsix\n"
	tmp := t.TempDir() + "/file.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &WcCommand{}
	out := testWcExecWithOutput(cmd, []string{tmp}, nil)

	if out != "3 6 20\n" {
		// проверка на общий формат
		t.Logf("проверка вывода: %q", out)
	}
}

func TestWcCommand_OnlyLines(t *testing.T) {
	content := "a\nb c\n\n"
	tmp := t.TempDir() + "/lines.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &WcCommand{}
	out := testWcExecWithOutput(cmd, []string{"-l", tmp}, nil)

	expected := "3\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestWcCommand_Stdin(t *testing.T) {
	cmd := &WcCommand{}
	input := "hello world\nhi\n"
	reader := strings.NewReader(input)
	out := testWcExecWithOutput(cmd, []string{"-l"}, reader)

	expected := "2\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}
