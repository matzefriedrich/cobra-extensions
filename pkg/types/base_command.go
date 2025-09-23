package types

import "context"

// BaseCommand A base type for typed commands.
type BaseCommand struct{}

// Execute performs the primary action defined by the BaseCommand.
func (c *BaseCommand) Execute(_ context.Context) {
}

var _ TypedCommand = (*BaseCommand)(nil)
