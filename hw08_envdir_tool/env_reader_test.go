package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testReadDir")
	if err != nil {
		t.Fatalf("Failed during create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	envFiles := []struct {
		name  string
		value string
	}{
		{"FOO", "value1"},
		{"EMPTY", ""},
		{"FOOBAR", "value\x00foobar   "},
	}
	os.Setenv("FOOBAR", "FOOBAR")

	for _, file := range envFiles {
		filePath := filepath.Join(tempDir, file.name)
		err = os.WriteFile(filePath, []byte(file.value), 0644)
		if err != nil {
			t.Fatalf("Failed during create file %s: %v", file.name, err)
		}
	}

	result, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("%v", err)
	}

	expected := Environment{"FOO": EnvValue{Value: "value1", NeedRemove: false}, "EMPTY": EnvValue{Value: "", NeedRemove: false}, "FOOBAR": EnvValue{Value: "value\nfoobar", NeedRemove: true}}
	require.Equal(t, expected, result)
	os.RemoveAll(tempDir)
}
