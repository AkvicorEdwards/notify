package parameter

// Command structure
type Arg struct {
	// Number of parameters required
	Size  int
	// All parameters that follow belong to this command
	// If the value is "true", "Size" will have no effect
	Block		bool
	// Execute "command"
	Executor func([]string)
}

// Command collection
type Args map[string]Arg

func GetArgs() Args {
	return Args{}
}
