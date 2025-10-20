package commands

type ExitCommand struct{}

func (e *ExitCommand) Name() string {
	return "exit"
}

func (e *ExitCommand) Exec(args []string) error {
	panic("not implemented")
}

func (e *ExitCommand) Help() string {
	panic("not implemented")
}
