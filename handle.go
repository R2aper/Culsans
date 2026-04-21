package main

import (
	"fmt"
	"os"
)

func handleAdd(name, keyPath, commit_message string, noCommit bool) {
	if keyPath == "" {
		fmt.Fprintf(os.Stderr, "Error:\n-k flag required for add command\n")
		os.Exit(1)
	}

	fmt.Printf("Enter content for '%s':\n", name)
	msg, err := readDataWithMask()
	defer func() {
		for i := range msg {
			msg[i] = 0
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError while reading:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	ciphertext, err := encryptWithPublicKey(keyPath, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting:\n%v\n", err)
		os.Exit(1)
	}

	filename := name + ".gpg"
	if err := os.WriteFile(filename, ciphertext, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' added\n", name)

	if !noCommit {
		// TODO: commit with message
	}
}

func handleRemove(name, commit_message string, noCommit bool) {
	filename := name + ".gpg"
	if err := os.Remove(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing password:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' removed\n", name)

	if !noCommit {
		// TODO: commit with message
	}
}

func handleShow(name, keyPath string) {
	if keyPath == "" {
		fmt.Fprintf(os.Stderr, "Error:\n-k flag required for show command\n")
		os.Exit(1)
	}

	fmt.Println("Enter passphrase:")
	pass, err := readDataWithMask()
	defer func() {
		for i := range pass {
			pass[i] = 0
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError while reading:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Println()

	filename := name + ".gpg"
	encryptedData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password file:\n%v\n", err)
		os.Exit(1)
	}

	data, err := decryptWithPrivateKey(keyPath, pass, encryptedData)
	defer func() {
		for i := range data {
			data[i] = 0
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decrypting:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Data:\n%s\n", string(data))
}
