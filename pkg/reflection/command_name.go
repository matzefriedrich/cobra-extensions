package reflection

// CommandName represents the name of a command in the system.
type CommandName struct {
	name string
}

// String returns the string representation of the CommandName instance.
func (c CommandName) String() string {
	return c.name
}
