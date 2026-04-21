// github.com/ProtonMail/gopenpgp/v3
// go-git

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const VERSION = "0.0"

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: cl [global-flags] <command> [args]
Password manager using PGP encryption

Global Flags:
  -k string
    	Path to PGP key file (public key for add, private key for show)
  -q	Don't commit changes to git
  -m string
    	Commit message
  -h	Show this help message
  -v	Show version

Commands:
  init		Initialize a new git repository
  list		List all passwords in the vault
  add <name>	Add a new password
  remove <name>	Remove a password
  show <name>	Show password content

Examples:
  cl -k ~/.pgp/pub.key add mypassword
  cl -k ~/.pgp/priv.key show mypassword
  cl list
`)
}

// Read data from stdout and mask input with '*'
// Stop reading after two newlines
func readDataWithMask() ([]byte, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return nil, fmt.Errorf("Stdin is not a terminal")
	}

	oldState, err := term.MakeRaw(fd)
	defer term.Restore(fd, oldState)
	if err != nil {
		return nil, fmt.Errorf("Error converting to raw mode: %v", err)
	}

	data := make([]byte, 0, 256)
	buf := make([]byte, 1)
	enterCount := 0

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("%v", err)
		}

		b := buf[0]

		// Enter
		if b == '\r' || b == '\n' {
			enterCount++
			fmt.Println()
			if enterCount == 2 {
				break
			}
			continue
		}

		// Backspace (\b) или Delete (\x7f)
		if b == '\b' || b == '\x7f' {
			if len(data) > 0 {
				data = data[:len(data)-1]
				fmt.Print("\b \b")
			}
			continue
		}

		// Ctrl+C
		if b == '\x03' {
			fmt.Println()
			return nil, fmt.Errorf("Ctrl+C pressed")
		}

		if b < 32 {
			continue
		}

		data = append(data, b)
		fmt.Print("*")
	}

	return data, nil
}

func main() {
	// Global flags
	help := flag.Bool("h", false, "Show help message")
	version := flag.Bool("v", false, "Show version")
	keyPath := flag.String("k", "", "Path to PGP key file")
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
		fmt.Println("TODO:Initializing password vault...")

	case "list":
		fmt.Println("TODO:Listing passwords...")

	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error:\nAdd requires a password name\n")
			os.Exit(1)
		}
		handleAdd(args[1], *keyPath, *message, *noCommit)

	case "remove":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error:\nRemove requires a password name\n")
			os.Exit(1)
		}
		handleRemove(args[1], *message, *noCommit)

	case "show":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error:\nShow requires a password name\n")
			os.Exit(1)
		}
		handleShow(args[1], *keyPath)

	default:
		fmt.Fprintf(os.Stderr, "Error:\nUnknown command '%s'\n", cmd)
		os.Exit(1)
	}
}
