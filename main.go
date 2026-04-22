package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const VERSION = "0.0"

func fatalError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: cl [global-flags] <command> [args]
Password manager using PGP encryption

Global Flags:
  -pub string		Path to public PGP key file public key for add
  -sec string 		Path to public PGP key file public key for show and signing commits
  -m string 		Commit message
  -q				Don't commit changes to git
  -s				Sign commit
  -h				Show this help message
  -v				Show version

Commands:
  init				Initialize a new git repository
  list				List all passwords in the vault
  add <name>		Add a new password
  remove <name>		Remove a password
  show <name>		Show password content

Examples:
  cl -k ~/.pgp/pub.key add mypassword
  cl -k ~/.pgp/priv.key show mypassword
  cl list
`)
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

func main() {
	// Global flags
	help := flag.Bool("h", false, "Show help message")
	version := flag.Bool("v", false, "Show version")
	pubKeyPath := flag.String("pub", "", "Path to public PGP key file")
	secKeyPath := flag.String("sec", "", "Path to private PGP key file")
	sign := flag.Bool("s", false, "Sign commit")
	noCommit := flag.Bool("q", false, "Don't commit changes")
	message := flag.String("m", "Update password vault", "Commit message")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		return
	}

	if *version {
		fmt.Fprintf(os.Stderr, "%s\n", VERSION)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	// Parsing commands
	cmd := args[0]
	switch cmd {
	case "init":
		handleInit()

	case "list":
		fmt.Println("TODO:Listing passwords...")

	case "add":
		if len(args) < 2 {
			fatalError("The 'add' command requires a password name as an argument")
		}

		var err error
		*pubKeyPath, err = parsePath(*pubKeyPath)
		if err != nil {
			fatalError("Invalid public key path: %v", err)
		}

		handleAdd(args[1], *pubKeyPath, *message, *noCommit, *sign, *secKeyPath)

	case "remove":
		if len(args) < 2 {
			fatalError("The 'remove' command requires a password name as an argument")
		}
		handleRemove(args[1], *message, *noCommit, *sign, *secKeyPath)

	case "show":
		if len(args) < 2 {
			fatalError("The 'show' command requires a password name as an argument")
		}

		var err error
		*secKeyPath, err = parsePath(*secKeyPath)
		if err != nil {
			fatalError("Invalid private key path: %v", err)
		}

		handleShow(args[1], *secKeyPath)

	default:
		fatalError("Unknown command '%s'", cmd)
	}
}
