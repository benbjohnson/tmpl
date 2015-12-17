package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// Extension is the required file extension for processed files.
const Extension = ".tmpl"

func main() {
	m := NewMain()
	if err := m.ParseFlags(os.Args[1:]); err != nil {
		fmt.Fprintln(m.Stderr, err)
		os.Exit(2)
	}

	if err := m.Run(); err != nil {
		fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	// Files to be processed.
	Paths []string

	// Data to be applied to the files during generation.
	Data interface{}

	// Standard input/output
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// ParseFlags parses the command line flags from args.
func (m *Main) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("tmp", flag.ContinueOnError)
	fs.SetOutput(m.Stderr)
	data := fs.String("data", "", "json data")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// Parse JSON data.
	if *data != "" {
		if err := json.Unmarshal([]byte(*data), &m.Data); err != nil {
			return err
		}
	}

	// All arguments are considered paths to process.
	m.Paths = fs.Args()

	return nil
}

// Run executes the program.
func (m *Main) Run() error {
	// Verify we have at least one path.
	if len(m.Paths) == 0 {
		return errors.New("path required")
	}

	// Process each path.
	for _, path := range m.Paths {
		if err := m.process(path); err != nil {
			return err
		}
	}

	return nil
}

// process reads a template file from path, processes it, and writes it to its generated path.
func (m *Main) process(path string) error {
	// Validate that we have a prefix we can strip off for the generated path.
	if !strings.HasSuffix(path, Extension) {
		return fmt.Errorf("path must have %s extension: %s", Extension, path)
	}

	// Stat the file to retrieve the mode.
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	} else if err != nil {
		return err
	}

	// Parse file into template.
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	// Execute template.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m.Data); err != nil {
		return err
	}

	// Write buffer to file.
	outputPath := strings.TrimSuffix(path, Extension)
	if err := ioutil.WriteFile(outputPath, buf.Bytes(), fi.Mode()); err != nil {
		return err
	}

	return nil
}
