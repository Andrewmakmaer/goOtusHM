package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("Copy non exist file", func(t *testing.T) {
		err := Copy("nonexist.txt", "nonexist.txt.copy", 0, 0)
		require.Truef(t, errors.Is(err, os.ErrNotExist), "actual error %q", err)
	})

	t.Run("offset is larger than file size", func(t *testing.T) {
		err := Copy("testdata/smallfile.txt", "small.txt.copy", 10, 0)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual error %q", err)
	})

	t.Run("offset is larger than file size", func(t *testing.T) {
		err := Copy("/dev/urandom", "small.txt.copy", 10, 0)
		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual error %q", err)
	})
}
