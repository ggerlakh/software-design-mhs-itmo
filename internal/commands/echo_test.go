package commands

import (
	"io"
	"os"
	"strings"
	"testing"
)

// Подменяет стандартный вывод, чтобы проверить результат Exec().
func captureEchoOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestEchoCommand_Simple(t *testing.T) {
	cmd := &EchoCommand{}

	out := captureEchoOutput(func() {
		_ = cmd.Exec([]string{"Hello", "world"})
	})

	expected := "Hello world\n"
	if out != expected {
		t.Errorf("ожидался вывод %q, получено %q", expected, out)
	}
}

func TestEchoCommand_NoNewline(t *testing.T) {
	cmd := &EchoCommand{}

	out := captureEchoOutput(func() {
		_ = cmd.Exec([]string{"-n", "Hello", "world"})
	})

	expected := "Hello world"
	if out != expected {
		t.Errorf("опция -n работает некорректно: ожидалось %q, получено %q", expected, out)
	}
}

func TestEchoCommand_EmptyArgs(t *testing.T) {
	cmd := &EchoCommand{}

	out := captureEchoOutput(func() {
		_ = cmd.Exec([]string{})
	})

	expected := "\n"
	if out != expected {
		t.Errorf("ожидался просто перевод строки при отсутствии аргументов: %q", out)
	}
}

func TestEchoCommand_OnlyFlagNoArgs(t *testing.T) {
	cmd := &EchoCommand{}

	out := captureEchoOutput(func() {
		_ = cmd.Exec([]string{"-n"})
	})

	expected := ""
	if out != expected {
		t.Errorf("при использовании -n без аргументов не должно быть вывода: %q", out)
	}
}

func TestEchoCommand_MultipleSpaces(t *testing.T) {
	cmd := &EchoCommand{}

	out := captureEchoOutput(func() {
		_ = cmd.Exec([]string{"foo", "bar", "baz"})
	})

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
