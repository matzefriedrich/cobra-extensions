package commands

import "github.com/matzefriedrich/cobra-extensions/pkg/types"

// BaseCommand A base type for typed commands.
type BaseCommand struct{}

// Execute performs the primary action defined by the BaseCommand.
func (c *BaseCommand) Execute() {
}

var _ types.TypedCommand = (*BaseCommand)(nil)
