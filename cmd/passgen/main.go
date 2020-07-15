package main

import (
	"github.com/spf13/cobra"
)

var (
	// This is used for platform-specific bounds checking of parsed uints.
	uintMax = ^uint(0)
)

func main() {
	// Define the root command.
	rootCmd := &cobra.Command{
		Use:   "passgen",
		Short: "Generate passwords and passphrases",
	}

	// Construct the password generation subcommand.
	passwordCmd := buildPasswordCmd()
	rootCmd.AddCommand(passwordCmd)

	// Run the root command.
	rootCmd.Execute()
}
