package executor

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/cli/internal/commands"
)

type mockBuiltin struct {
	name   string
	args   []string
	called bool
}

func (m *mockBuiltin) Name() string {
	return m.name
}

func (m *mockBuiltin) Exec(args []string, ctx *commands.CommandContext) error {
	m.called = true
	m.args = append([]string{}, args...)
	return nil
}

func (m *mockBuiltin) Help() string {
	return ""
}

type funcBuiltin struct {
	name string
	run  func(args []string, ctx *commands.CommandContext) error
}

func (f *funcBuiltin) Name() string { return f.name }
func (f *funcBuiltin) Help() string { return "" }
func (f *funcBuiltin) Exec(args []string, ctx *commands.CommandContext) error {
	if f.run != nil {
		return f.run(args, ctx)
	}
	return nil
}

func TestExecutor_ExecuteBuiltin(t *testing.T) {
	builtin := &mockBuiltin{name: "mock"}
	ex := NewExecutor(map[string]string{}, []commands.BuiltinCommand{builtin})

	ex.Execute(Plan{
		Commands: []ExecutableCommand{
			{Name: "mock", Args: []string{"hello", "world"}},
		},
	})

	if !builtin.called {
		t.Fatalf("ожидалось выполнение builtin команды")
	}

	if len(builtin.args) != 2 || builtin.args[0] != "hello" {
		t.Fatalf("аргументы переданы некорректно: %#v", builtin.args)
	}
}

func TestExecutor_EnvAssignment(t *testing.T) {
	env := map[string]string{}
	ex := NewExecutor(env, nil)

	ex.Execute(Plan{
		Commands: []ExecutableCommand{
			{Name: "FOO=bar"},
		},
	})

	if env["FOO"] != "bar" {
		t.Fatalf("переменная окружения не установлена")
	}
}

func TestExecutor_ExecutePipelineBuiltinFlow(t *testing.T) {
	producer := &funcBuiltin{
		name: "produce",
		run: func(args []string, ctx *commands.CommandContext) error {
			_, err := ctx.Stdout.Write([]byte("data"))
			return err
		},
	}

	consumer := &funcBuiltin{
		name: "consume",
		run: func(args []string, ctx *commands.CommandContext) error {
			payload, err := io.ReadAll(ctx.Stdin)
			if err != nil {
				return err
			}
			_, err = ctx.Stdout.Write([]byte("processed:" + string(payload)))
			return err
		},
	}

	ex := NewExecutor(map[string]string{}, []commands.BuiltinCommand{producer, consumer})

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = oldStdout
		_ = r.Close()
	}()

	ex.Execute(Plan{
		Commands: []ExecutableCommand{
			{Name: "produce"},
			{Name: "consume"},
		},
	})

	_ = w.Close()
	output, _ := io.ReadAll(r)
	if string(bytes.TrimSpace(output)) != "processed:data" {
		t.Fatalf("ожидался вывод processed:data, получено: %q", string(output))
	}
}

func TestExecutor_ExecuteExternalCommand(t *testing.T) {
	ex := NewExecutor(map[string]string{}, nil)

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = oldStdout
		_ = r.Close()
	}()

	ex.Execute(Plan{
		Commands: []ExecutableCommand{
			{Name: "printf", Args: []string{"external-output"}},
		},
	})

	_ = w.Close()
	output, _ := io.ReadAll(r)
	if !strings.Contains(string(output), "external-output") {
		t.Fatalf("внешняя команда не записала ожидаемый вывод: %q", string(output))
	}
}
