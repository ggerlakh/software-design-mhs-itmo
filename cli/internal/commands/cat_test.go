package commands

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// testExecWithOutput выполняет команду с переданными аргументами и возвращает вывод
func testExecWithOutput(cmd CommandExecutor, args []string, stdin io.Reader) string {
	var buf bytes.Buffer
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

func TestCatCommand_SimpleFile(t *testing.T) {
	content := "Hello\nWorld\n"
	tmp := t.TempDir() + "/test.txt"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &CatCommand{}
	out := testExecWithOutput(cmd, []string{tmp}, nil)

	if out != content {
		t.Errorf("ожидался вывод %q, получено %q", content, out)
	}
}

func TestCatCommand_NumberNonEmpty(t *testing.T) {
	content := "a\n\nb\n"
	tmp := t.TempDir() + "/nonempty.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &CatCommand{}
	out := testExecWithOutput(cmd, []string{"-b", tmp}, nil)

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if !strings.Contains(lines[0], "1") || !strings.Contains(lines[2], "2") {
		t.Errorf("ожидалось нумерование только непустых строк:\n%s", out)
	}
}

func TestCatCommand_SqueezeBlank(t *testing.T) {
	content := "1\n\n\n2\n\n\n\n3\n"
	tmp := t.TempDir() + "/squeeze.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &CatCommand{}
	out := testExecWithOutput(cmd, []string{"-s", tmp}, nil)

	expected := "1\n\n2\n\n3\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestCatCommand_ShowEndsAndTabs(t *testing.T) {
	content := "a\tb\n"
	tmp := t.TempDir() + "/tabs.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &CatCommand{}
	out := testExecWithOutput(cmd, []string{"-E", "-T", tmp}, nil)

	expected := "a^Ib$\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestCatCommand_Stdin(t *testing.T) {
	cmd := &CatCommand{}
	input := "input from stdin\n"
	reader := strings.NewReader(input)
	out := testExecWithOutput(cmd, []string{"-"}, reader)

	if out != "input from stdin\n" {
		t.Errorf("ожидался вывод stdin, получено %q", out)
	}
}
