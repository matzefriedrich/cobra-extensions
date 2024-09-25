package types

// CommandArgs stores argument options for a command.
type CommandArgs struct {
	MinimumArgs int
}

// CommandArgsOption defines a function type that modifies a CommandArgs instance.
type CommandArgsOption func(*CommandArgs)

// MinimumArgumentsRequired sets the minimum number of arguments required for a command.
func MinimumArgumentsRequired(value int) CommandArgsOption {
	return func(args *CommandArgs) {
		args.MinimumArgs = value
	}
}

// NewCommandArgs creates a new CommandArgs instance, applying given CommandArgsOption functions.
func NewCommandArgs(options ...CommandArgsOption) CommandArgs {
	args := CommandArgs{}
	for _, option := range options {
		option(&args)
	}
	return args
}
