package main

import (
	"fmt"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/go-git/go-git/v5"
)

func handleInit() {
	_, err := git.PlainInit("./", false)
	if err != nil {
		fatalError("Failed to initialize Git repository: %v", err)
	}

	dir, _ := os.Getwd()
	fmt.Println("Initialized git repository at", dir)
}

func handleAdd(name string, pubKeyPath string, commit_message string, noCommit bool, sign bool, secKeyPath string) {
	if pubKeyPath == "" {
		fatalError("-pub flag is required for the 'add' command")
	}

	if secKeyPath == "" && sign {
		fatalError("-sec flag is required when signing commits")
	}

	fmt.Printf("Enter content for '%s':\n", name)
	msg, err := readDataWithMask(true)
	defer func() {
		for i := range msg {
			msg[i] = 0
		}
	}()
	if err != nil {
		fatalError("Failed to read input: %v", err)
	}

	ciphertext, err := encryptWithPublicKey(pubKeyPath, msg)
	if err != nil {
		fatalError("Failed to encrypt data: %v", err)
	}

	filename := name + ".gpg"
	if err := os.WriteFile(filename, ciphertext, 0644); err != nil {
		fatalError("Failed to write encrypted file: %v", err)
	}

	if !noCommit {
		// Open repo
		rep, err := git.PlainOpen("./")
		if err != nil {
			fatalError("Failed to open Git repository: %v", err)
		}

		var openpgpEntity *openpgp.Entity

		if sign {
			fmt.Println("Enter passphrase:")
			pass, err := readDataWithMask(false)
			defer func() {
				for i := range pass {
					pass[i] = 0
				}
			}()
			if err != nil {
				fatalError("Failed to read passphrase: %v", err)
			}
			fmt.Println()

			openpgpEntity, err = getSigningEntity(secKeyPath, pass)
		} else {
			openpgpEntity = nil
		}

		_, err = CommitChanges(rep, []string{filename}, commit_message, openpgpEntity)
		if err != nil {
			fatalError("Failed to commit changes: %v", err)
		}
	}

	fmt.Printf("Password '%s' added\n", name)
}

func handleRemove(name string, commit_message string, noCommit bool, sign bool, secKeyPath string) {
	filename := name + ".gpg"
	if err := os.Remove(filename); err != nil {
		fatalError("Failed to remove password file: %v", err)
	}

	if !noCommit {
		// Open repo
		rep, err := git.PlainOpen("./")
		if err != nil {
			fatalError("Failed to open Git repository: %v", err)
		}

		var openpgpEntity *openpgp.Entity

		if sign {
			fmt.Println("Enter passphrase:")
			pass, err := readDataWithMask(false)
			defer func() {
				for i := range pass {
					pass[i] = 0
				}
			}()
			if err != nil {
				fatalError("Failed to read passphrase: %v", err)
			}
			fmt.Println()

			openpgpEntity, err = getSigningEntity(secKeyPath, pass)
		} else {
			openpgpEntity = nil
		}

		_, err = CommitChanges(rep, []string{filename}, commit_message, openpgpEntity)
		if err != nil {
			fatalError("Failed to commit changes: %v", err)
		}
	}

	fmt.Printf("Password '%s' removed\n", name)
}

func handleShow(name string, keyPath string) {
	if keyPath == "" {
		fatalError("-sec flag is required for the 'show' command")
	}

	fmt.Println("Enter passphrase:")
	pass, err := readDataWithMask(false)
	defer func() {
		for i := range pass {
			pass[i] = 0
		}
	}()
	if err != nil {
		fatalError("Failed to read passphrase: %v", err)
	}
	fmt.Println()

	filename := name + ".gpg"
	encryptedData, err := os.ReadFile(filename)
	if err != nil {
		fatalError("Failed to read password file: %v", err)
	}

	data, err := decryptWithPrivateKey(keyPath, pass, encryptedData)
	defer func() {
		for i := range data {
			data[i] = 0
		}
	}()
	if err != nil {
		fatalError("Failed to decrypt data: %v", err)
	}

	fmt.Printf("Content of %s:\n%s\n", name, string(data))
}
