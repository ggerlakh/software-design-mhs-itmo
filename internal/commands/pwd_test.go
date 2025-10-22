package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// testExecWithOutput выполняет команду с переданными аргументами и возвращает вывод
func testPwdExecWithOutput(cmd CommandExecutor, args []string) string {
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

func TestPwdCommand_Simple(t *testing.T) {
	cmd := &PwdCommand{}

	expected, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	out := testPwdExecWithOutput(cmd, nil)

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

	out := testPwdExecWithOutput(cmd, []string{"extra", "args"})

	out = strings.TrimSpace(out)
	if out != expected {
		t.Errorf("pwd не должен зависеть от аргументов: ожидалось %q, получено %q", expected, out)
	}
}
