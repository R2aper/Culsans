// github.com/ProtonMail/gopenpgp/v3
// go-git

package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"
)

func usage() {
	fmt.Println("" +
		"Usage: cl [Options] <command> \n" +
		"Password manager\n" +
		"\nOptions:" +
		"-h\tShow this help message\n" +
		"-v\t\tShow version\n" +
		"-q\t\tDon't commit changes\n" +
		"-m\tSpecify commit message\n" +
		"\nCommands:\n" +
		"init\t\t\tInitialize a new git repository in the current working directory(Similar to git init)\n" +
		"list\t\t\tList all passwords in the vault\n" +
		"add <name>\t\tAdd a new password\n" +
		"remove <name>\t\tRemove a password\n" +
		"show <name>\t\tShow content of password")
}

func main() {
	// Flags
	help_flag := flag.Bool("h", false, "Print usage")
	version_flag := flag.Bool("v", false, "Print version")

	// States
	key_value := flag.String("k", "", "Specify public/private key path")

	// Modes
	add_flag := flag.String("add", "", "Specify output file name")
	remove_flag := flag.String("remove", "", "TODO")
	show_flag := flag.String("show", "", "Specify input file name")

	flag.Parse()

	if *help_flag {
		usage()
		return
	}

	if *version_flag {
		fmt.Println("0.0")
		return
	}

	if *add_flag != "" {
		fmt.Println("Enter message:")
		msg, err := term.ReadPassword(int(os.Stdin.Fd()))
		defer func() {
			for i := range msg {
				msg[i] = 0
			}
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError while reading: %v\n", err)
			return
		}
		fmt.Println()

		ciphertext, err := encryptWithPublicKey(*key_value, msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		os.WriteFile(*add_flag, ciphertext, 0644)
		// TODO: commit + sign
		return
	}

	if *remove_flag != "" {
		fmt.Println("TODO")
	}

	if *show_flag != "" {
		fmt.Println("Enter passphrase:")
		pass, err := term.ReadPassword(int(os.Stdin.Fd()))
		defer func() {
			for i := range pass {
				pass[i] = 0
			}
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError while reading: %v\n", err)
			return
		}
		fmt.Println()

		msg, err := os.ReadFile(*show_flag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		data, err := decryptWithPrivateKey(*key_value, pass, msg)
		defer func() {
			for i := range data {
				data[i] = 0
			}
		}()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Printf("Data:\n%s\n", string(data))
		return
	}

}
