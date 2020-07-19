package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// This is used for platform-specific bounds checking of parsed uints.
	uintMax = ^uint(0)

	// Exit function used to aid test coverage.
	exitFunc = os.Exit

	// This version variable is populated at compilation.
	version string
)

func main() {
	// Define the root command.
	rootCmd := &cobra.Command{
		Use:     "passgen",
		Short:   "Generate passwords and passphrases",
		Version: version,
	}

	// Construct the password generation subcommand.
	passwordCmd := buildPasswordCmd()
	rootCmd.AddCommand(passwordCmd)

	// Construct the passphrase generation subcommand.
	passphraseCmd := buildPassphraseCmd()
	rootCmd.AddCommand(passphraseCmd)

	// Run the root command.
	err := rootCmd.Execute()
	if err != nil {
		exitFunc(2)
	}
}
