package commands

// cryptCommand A base command type for crypt commands.
type cryptCommand struct {
	Passphrase string `flag:"passphrase" usage:"The passphrase protecting a message."`
}
