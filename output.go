package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

// MultiWriter wraps multiple writers to write to both stdout and a file
type OutputWriter struct {
	file   *os.File
	stdout io.Writer
}

// NewOutputWriter creates a new output writer that writes to both stdout and result.txt
func NewOutputWriter() (*OutputWriter, error) {
	file, err := os.Create("result.txt")
	if err != nil {
		return nil, err
	}
	return &OutputWriter{
		file:   file,
		stdout: os.Stdout,
	}, nil
}

// Write writes to both stdout and file
func (w *OutputWriter) Write(p []byte) (n int, err error) {
	// Write to stdout
	fmt.Fprint(w.stdout, string(p))
	// Write to file
	return w.file.Write(p)
}

// Close closes the output file
func (w *OutputWriter) Close() error {
	return w.file.Close()
}

// Printf formats and writes to both outputs
func (w *OutputWriter) Printf(format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

// Println writes a line to both outputs
func (w *OutputWriter) Println(a ...interface{}) {
	fmt.Fprintln(w, a...)
}

// GetTimestamp returns current time in HH:MM:SS format
func GetTimestamp() string {
	return time.Now().Format("15:04:05")
}
