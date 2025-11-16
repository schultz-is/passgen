package main

import (
	"fmt"
	"strconv"
)

// parseUintArg parses a uint argument from command line args at the specified index.
func parseUintArg(args []string, index int, paramName string) (uint, error) {
	if index >= len(args) {
		return 0, nil // Not provided, use default
	}

	val, err := strconv.ParseUint(args[index], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s provided", paramName)
	}

	if val > uint64(uintMax) {
		return 0, fmt.Errorf("invalid %s provided", paramName)
	}

	return uint(val), nil
}
