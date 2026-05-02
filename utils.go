package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"golang.org/x/term"
)

func fatalError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

// Read data from stdout and mask input with '*'
// If endWithNewLine is true then reading stops after a newline that follows an empty line
func readDataWithMask(endWithNewLine bool) ([]byte, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil, fmt.Errorf("Stdin is not a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("Error converting to raw mode: %v", err)
	}
	defer term.Restore(fd, oldState)

	data := make([]byte, 0, 512)
	buf := make([]byte, 1)
	lineStart := 0

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		b := buf[0]

		// === Enter ===
		if b == '\r' || b == '\n' {
			if endWithNewLine {
				data = append(data, '\n')
				break
			}

			// === Multi-line mode ===
			if len(data) == lineStart {
				// Emtpy line
				fmt.Print("\n")
				break
			}

			// Have content
			data = append(data, '\n')
			fmt.Print("\n")
			lineStart = len(data)
			continue
		}

		// === Backspace ===
		if b == '\b' || b == '\x7f' {
			if len(data) > lineStart {
				data = data[:len(data)-1]
				fmt.Print("\b \b")
			}
			continue
		}

		// === Ctrl+C ===
		if b == '\x03' {
			fmt.Print("\n")
			return nil, fmt.Errorf("Ctrl+C pressed")
		}

		if b < 32 && b != '\t' {
			continue
		}

		data = append(data, b)
		fmt.Print("*")
	}

	data = data[:len(data)-1] // Remove \n
	return data, nil
}

func parsePath(p string) (string, error) {
	new_path := ""
	if len(p) >= 2 && p[0:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		new_path = path.Join(home, p[2:])
		return new_path, nil
	}

	return p, nil
}
