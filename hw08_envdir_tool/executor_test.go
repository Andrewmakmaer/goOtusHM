package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		cmd := []string{"echo", "hello"}
		env := Environment{}
		result := RunCmd(cmd, env)
		require.Equal(t, 0, result)
	})

	t.Run("non zero code", func(t *testing.T) {
		cmd := []string{"false"}
		env := Environment{}
		result := RunCmd(cmd, env)
		require.Equal(t, 1, result)
	})
}
