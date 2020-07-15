package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainCommand(t *testing.T) {
	require.NotPanics(t, main)
}
