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
		msg, err := readDataWithMask()
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
		pass, err := readDataWithMask()
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
