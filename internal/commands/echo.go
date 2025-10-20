package commands

type EchoCommand struct{}

func (e *EchoCommand) Name() string {
	return "echo"
}

func (e *EchoCommand) Exec(args []string) error {
	panic("not implemented")
}

func (e *EchoCommand) Help() string {
	panic("not implemented")
}
