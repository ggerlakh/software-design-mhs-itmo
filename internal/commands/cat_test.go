package commands

import (
	"io"
	"os"
	"strings"
	"testing"
)

// Подменяет стандартный вывод, чтобы проверить результат Exec().
func captureCatOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestCatCommand_SimpleFile(t *testing.T) {
	content := "Hello\nWorld\n"
	tmp := t.TempDir() + "/test.txt"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := &CatCommand{}
	out := captureCatOutput(func() {
		_ = cmd.Exec([]string{tmp})
	})

	if out != content {
		t.Errorf("ожидался вывод %q, получено %q", content, out)
	}
}

func TestCatCommand_NumberNonEmpty(t *testing.T) {
	content := "a\n\nb\n"
	tmp := t.TempDir() + "/nonempty.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &CatCommand{}
	out := captureCatOutput(func() {
		_ = cmd.Exec([]string{"-b", tmp})
	})

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
	out := captureCatOutput(func() {
		_ = cmd.Exec([]string{"-s", tmp})
	})

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
	out := captureCatOutput(func() {
		_ = cmd.Exec([]string{"-E", "-T", tmp})
	})

	expected := "a^Ib$\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestCatCommand_Stdin(t *testing.T) {
	cmd := &CatCommand{}
	r, w, _ := os.Pipe()
	os.Stdin = r
	io.WriteString(w, "input from stdin\n")
	w.Close()

	out := captureCatOutput(func() {
		_ = cmd.Exec([]string{"-"})
	})

	if out != "input from stdin\n" {
		t.Errorf("ожидался вывод stdin, получено %q", out)
	}
}
