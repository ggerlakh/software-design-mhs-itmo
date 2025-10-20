package commands

type WcCommand struct{}

func (w *WcCommand) Name() string {
	return "pwd"
}

func (w *WcCommand) Exec(args []string) error {
	panic("not implemented")
}

func (w *WcCommand) Help() string {
	panic("not implemented")
}
