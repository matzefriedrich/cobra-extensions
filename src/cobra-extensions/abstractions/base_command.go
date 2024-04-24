package abstractions

type ExecutableCommand interface {
	Execute()
}

// BaseCommand A base type for typed commands.
type BaseCommand struct{}

func (c *BaseCommand) Execute() {
}
