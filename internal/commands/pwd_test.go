package commands

import (
	"io"
	"os"
	"strings"
	"testing"
)

// Подменяет стандартный вывод, чтобы проверить результат Exec().
func capturePwdOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestPwdCommand_Simple(t *testing.T) {
	cmd := &PwdCommand{}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	out := capturePwdOutput(func() {
		_ = cmd.Exec(nil)
	})

	out = strings.TrimSpace(out)
	if out != expected {
		t.Errorf("ожидался текущий каталог %q, получено %q", expected, out)
	}
}

func TestPwdCommand_IgnoresArgs(t *testing.T) {
	cmd := &PwdCommand{}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	out := capturePwdOutput(func() {
		_ = cmd.Exec([]string{"extra", "args"})
	})

	out = strings.TrimSpace(out)
	if out != expected {
		t.Errorf("pwd не должен зависеть от аргументов: ожидалось %q, получено %q", expected, out)
	}
}
