package abstractions

type CommandArgs struct {
	MinimumArgs int
}

func NewCommandArgs(minimumArgumentsRequired int) CommandArgs {
	return CommandArgs{
		MinimumArgs: minimumArgumentsRequired,
	}
}
