package ape

type Command struct {
	name string
	args []string
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Args() []string {
	return c.args
}

func newCommand(name string, args []string) *Command {
	return &Command{
		name: name,
		args: args,
	}
}
