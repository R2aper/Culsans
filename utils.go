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

	data := make([]byte, 0, 256)
	buf := make([]byte, 1)
	lineStartIdx := 0 // Index in data where the current line starts

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("%v", err)
		}

		b := buf[0]

		// Handle Enter (\r or \n)
		if b == '\r' || b == '\n' {
			if !endWithNewLine {
				break
			}

			// If current line has no content (empty line), stop after this newline
			if len(data) == lineStartIdx {
				data = append(data, b)
				fmt.Println()
				break
			}
			// Line has content - continue reading
			data = append(data, b)
			fmt.Println()
			lineStartIdx = len(data) // Next line starts after this newline
			continue
		}

		// Handle Backspace (\b) or Delete (\x7f)
		if b == '\b' || b == '\x7f' {
			if len(data) > lineStartIdx {
				data = data[:len(data)-1]
				fmt.Print("\b \b")
			}
			continue
		}

		// Handle Ctrl+C
		if b == '\x03' {
			fmt.Println()
			return nil, fmt.Errorf("Ctrl+C pressed")
		}

		// Skip other control characters
		if b < 32 {
			continue
		}

		// Add printable character
		data = append(data, b)
		fmt.Print("*")
	}

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
