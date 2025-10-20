package commands

type CommandExecutor interface {
	Exec(args []string) error
}

type BuiltinCommand interface {
	CommandExecutor
	Name() string
	Help() string
}
