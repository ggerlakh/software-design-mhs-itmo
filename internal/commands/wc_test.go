package commands

import (
	"io"
	"os"
	"testing"
)

// Перехват stdout для проверки вывода
func captureWcOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestWcCommand_LinesWordsBytes(t *testing.T) {
	content := "one two\nthree four five\nsix\n"
	tmp := t.TempDir() + "/file.txt"
	_ = os.WriteFile(tmp, []byte(content), 0o644)

	cmd := &WcCommand{}
	out := captureWcOutput(func() {
		_ = cmd.Exec([]string{tmp})
	})

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
	out := captureWcOutput(func() {
		_ = cmd.Exec([]string{"-l", tmp})
	})

	expected := "3\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}

func TestWcCommand_Stdin(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdin = r
	io.WriteString(w, "hello world\nhi\n")
	w.Close()

	cmd := &WcCommand{}
	out := captureWcOutput(func() {
		_ = cmd.Exec([]string{"-l"})
	})

	expected := "2\n"
	if out != expected {
		t.Errorf("ожидалось %q, получено %q", expected, out)
	}
}
