package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// testExecWithOutput выполняет команду с переданными аргументами и возвращает вывод
func testEchoExecWithOutput(cmd CommandExecutor, args []string) string {
	var buf bytes.Buffer
	ctx := &CommandContext{
		Stdin:  os.Stdin,
		Stdout: &buf,
		Stderr: os.Stderr,
		Env:    make(map[string]string),
		Dir:    ".",
	}
	cmd.Exec(args, ctx)
	return buf.String()
}

func TestEchoCommand_Simple(t *testing.T) {
	cmd := &EchoCommand{}

	out := testEchoExecWithOutput(cmd, []string{"Hello", "world"})

	expected := "Hello world\n"
	if out != expected {
		t.Errorf("ожидался вывод %q, получено %q", expected, out)
	}
}

func TestEchoCommand_NoNewline(t *testing.T) {
	cmd := &EchoCommand{}

	out := testEchoExecWithOutput(cmd, []string{"-n", "Hello", "world"})

	expected := "Hello world"
	if out != expected {
		t.Errorf("опция -n работает некорректно: ожидалось %q, получено %q", expected, out)
	}
}

func TestEchoCommand_EmptyArgs(t *testing.T) {
	cmd := &EchoCommand{}

	out := testEchoExecWithOutput(cmd, []string{})

	expected := "\n"
	if out != expected {
		t.Errorf("ожидался просто перевод строки при отсутствии аргументов: %q", out)
	}
}

func TestEchoCommand_OnlyFlagNoArgs(t *testing.T) {
	cmd := &EchoCommand{}

	out := testEchoExecWithOutput(cmd, []string{"-n"})

	expected := ""
	if out != expected {
		t.Errorf("при использовании -n без аргументов не должно быть вывода: %q", out)
	}
}

func TestEchoCommand_MultipleSpaces(t *testing.T) {
	cmd := &EchoCommand{}

	out := testEchoExecWithOutput(cmd, []string{"foo", "bar", "baz"})

	expectedWords := []string{"foo", "bar", "baz"}
	for _, w := range expectedWords {
		if !strings.Contains(out, w) {
			t.Errorf("в выводе отсутствует слово %q: %q", w, out)
		}
	}

	if !strings.HasSuffix(out, "\n") {
		t.Errorf("ожидался перевод строки в конце вывода: %q", out)
	}
}
