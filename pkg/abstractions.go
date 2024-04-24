package pkg

type ExecutableCommand interface {
	Execute()
}

type BaseCommand struct{}

func (c *BaseCommand) Execute() {
}
