package abstractions

// ExecutableCommand defines an interface for commands that can be executed.
type ExecutableCommand interface {

	// Execute runs the command, performing its specific action based on the implementation.
	Execute()
}

// BaseCommand A base type for typed commands.
type BaseCommand struct{}

// Execute performs the primary action defined by the BaseCommand.
func (c *BaseCommand) Execute() {
}
