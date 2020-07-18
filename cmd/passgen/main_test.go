package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainCommand(t *testing.T) {
	// Ensure the main function doesn't panic with default arguments.
	require.NotPanics(t, main)

	// Set exit function to update local variable.
	var exitCode int
	exitFunc = func(code int) {
		exitCode = code
	}

	// Back up and override command arguments.
	args := make([]string, len(os.Args))
	_ = copy(args, os.Args)
	os.Args = append(os.Args, "x")

	// Run the main function with modified exit function and arguments.
	main()

	// Ensure we exited as expected.
	require.Equal(t, 2, exitCode)

	// Reset exit function and arguments.
	exitFunc = os.Exit
	os.Args = args
}
