package cli

// CLICommon holds common stuff
type CLICommon struct {
	// Debug will enable debug mode
	Debug    DebugFlag    `help:"Enable debug mode." short:"v" env:"DEBUG" default:"false"`
	Terminal TerminalFlag `help:"Enable terminal mode ( log not printed as json)" env:"TERMINAL" default:"false"`
}

// CLI is the overall cli-struct
var CLI struct {
	CLICommon
}
