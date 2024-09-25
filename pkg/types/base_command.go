package types

// BaseCommand A base type for typed commands.
type BaseCommand struct{}

// Execute performs the primary action defined by the BaseCommand.
func (c *BaseCommand) Execute() {
}

var _ TypedCommand = (*BaseCommand)(nil)
