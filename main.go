package main

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "0.0"

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
		handleList()

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
