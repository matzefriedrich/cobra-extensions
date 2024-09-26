package utils

import (
	"bytes"
	"os"
	"testing"
)

// CaptureStdout captures the standard output of the provided callback function.
//
// Parameters:
//   - t: the testing framework instance.
//   - in: the input string to be written to stdin.
//   - cb: the callback function whose stdout needs to be captured.
//
// Returns:
//   - The captured stdout as a string.
func CaptureStdout(t *testing.T, in string, cb func() error) (string, error) {

	originalStdin := os.Stdin
	defer func() { os.Stdin = originalStdin }()

	originalStdOut := os.Stdout
	defer func() { os.Stdout = originalStdOut }()

	stdinReader, stdinWriter, stdinPipeErr := os.Pipe()
	if stdinPipeErr != nil {
		t.Fatalf("failed to create pipe: %v", stdinPipeErr)
	}

	stdoutReader, stdoutWriter, stdoutPipeErr := os.Pipe()
	if stdoutPipeErr != nil {
		t.Fatalf("failed to create pipe: %v", stdoutPipeErr)
	}

	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	_, _ = stdinWriter.Write([]byte(in))
	_ = stdinWriter.Close()

	output := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(stdoutReader)
		output <- buf.String()
	}()

	err := cb()
	if err != nil {
		return "", err
	}

	_ = stdoutWriter.Close()

	data := <-output
	return data, nil
}
