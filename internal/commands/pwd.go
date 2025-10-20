package commands

type PwdCommand struct{}

func (p *PwdCommand) Name() string {
	return "pwd"
}

func (p *PwdCommand) Exec(args []string) error {
	panic("not implemented")
}

func (p *PwdCommand) Help() string {
	panic("not implemented")
}
