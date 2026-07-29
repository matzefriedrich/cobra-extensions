package types

// CommandName represents the name of a command in the system.
//
// Deprecated: Use BaseCommand and cobra-x tags instead.
type CommandName struct {
	name string
}

// String returns the string representation of the CommandName instance.
func (c CommandName) String() string {
	return c.name
}
