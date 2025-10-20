package commands

type CatCommand struct{}

func (c *CatCommand) Name() string {
	return "cat"
}

func (c *CatCommand) Exec(args []string) error {
	panic("not implemented")
}

func (c *CatCommand) Help() string {
	panic("not implemented")
}
