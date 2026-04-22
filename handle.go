package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func handleInit() {
	_, err := git.PlainInit("./", false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Vault initialized")
}

func handleAdd(name string, keyPath string, commit_message string, noCommit bool) {
	if keyPath == "" {
		fmt.Fprintf(os.Stderr, "Error:\n-pub flag required for add command\n")
		os.Exit(1)
	}

	fmt.Printf("Enter content for '%s':\n", name)
	msg, err := readDataWithMask(true)
	defer func() {
		for i := range msg {
			msg[i] = 0
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError while reading:\n%v\n", err)
		os.Exit(1)
	}

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
		// Open repo
		repo, err := git.PlainOpen("./")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		w, err := repo.Worktree()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		// Stage file
		_, err = w.Add(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		// Get author signature
		sig, err := GetASignature(repo)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Commit(commit_message, &git.CommitOptions{Author: sig})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}
	}
}

func handleRemove(name string, commit_message string, noCommit bool) {
	filename := name + ".gpg"
	if err := os.Remove(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing password:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' removed\n", name)

	if !noCommit {
		// Open repo
		repo, err := git.PlainOpen("./")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		w, err := repo.Worktree()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		// Stage file
		_, err = w.Add(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}

		// Get author signature
		sig, err := GetASignature(repo)
		if err != nil {
			log.Fatal(err)
		}

		_, err = w.Commit(commit_message, &git.CommitOptions{Author: sig})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}
	}
}

func handleShow(name string, keyPath string) {
	if keyPath == "" {
		fmt.Fprintf(os.Stderr, "Error:\n-sec flag required for show command\n")
		os.Exit(1)
	}

	fmt.Println("Enter passphrase:")
	pass, err := readDataWithMask(false)
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
