package interpreter

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/ggerlakh/software-design-mhs-itmo/internal/commands"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/executor"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/parser"
	"github.com/ggerlakh/software-design-mhs-itmo/internal/preprocessor"
)

func TestInterpreter_Creation(t *testing.T) {
	i := &Interpreter{
		Preprocessor: preprocessor.NewPreprocessor(),
		Parser:       parser.NewParser([]string{"echo"}),
		Executor:     executor.NewExecutor(map[string]string{}, nil),
	}

	if i.Preprocessor == nil || i.Parser == nil || i.Executor == nil {
		t.Fatalf("интерпретатор должен содержать все слои")
	}
}

func TestToExecutionPlan(t *testing.T) {
	p := parser.Pipeline{
		Commands: []parser.ParsedCommand{
			{Name: "echo", Args: []string{"hello"}},
			{Name: "wc", Args: []string{"-l"}},
		},
	}

	plan := toExecutionPlan(p)

	if len(plan.Commands) != 2 {
		t.Fatalf("ожидалось 2 команды, получено: %d", len(plan.Commands))
	}

	if plan.Commands[0].Name != "echo" || plan.Commands[1].Name != "wc" {
		t.Fatalf("команды сконвертированы неверно: %#v", plan.Commands)
	}
}

type testBuiltin struct {
	name string
	run  func(args []string, ctx *commands.CommandContext) error
}

func (tb *testBuiltin) Name() string { return tb.name }
func (tb *testBuiltin) Help() string { return "" }
func (tb *testBuiltin) Exec(args []string, ctx *commands.CommandContext) error {
	if tb.run != nil {
		return tb.run(args, ctx)
	}
	return nil
}

func TestInterpreter_StartRunsCommands(t *testing.T) {
	env := map[string]string{"TARGET": "world"}
	pre := preprocessor.NewPreprocessor(&preprocessor.EnvSubstitutionStep{Env: env})
	par := parser.NewParser([]string{"greet"})

	var executed bool
	greet := &testBuiltin{
		name: "greet",
		run: func(args []string, ctx *commands.CommandContext) error {
			executed = true
			_, err := ctx.Stdout.Write([]byte("hello " + strings.Join(args, " ")))
			return err
		},
	}
	exec := executor.NewExecutor(env, []commands.BuiltinCommand{greet})

	inputReader, inputWriter, _ := os.Pipe()
	_, _ = inputWriter.WriteString("greet $TARGET\nexit\n")
	_ = inputWriter.Close()

	oldStdin := os.Stdin
	os.Stdin = inputReader

	outputReader, outputWriter, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = outputWriter

	interpreter := &Interpreter{
		Preprocessor: pre,
		Parser:       par,
		Executor:     exec,
	}
	interpreter.Start()

	_ = outputWriter.Close()
	output, _ := io.ReadAll(outputReader)

	os.Stdin = oldStdin
	os.Stdout = oldStdout

	if !executed {
		t.Fatalf("ожидалось выполнение команды greet")
	}

	if !strings.Contains(string(output), "hello world") {
		t.Fatalf("не найден ожидаемый вывод: %q", string(output))
	}
}
