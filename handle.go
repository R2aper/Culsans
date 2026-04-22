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
		fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Println("Vault initialized")
}

func handleAdd(name string, pubKeyPath string, commit_message string, noCommit bool, sign bool, secKeyPath string) {
	if pubKeyPath == "" {
		fmt.Fprintf(os.Stderr, "Error:\n-pub flag required for add command\n")
		os.Exit(1)
	}

	if secKeyPath == "" && sign {
		fmt.Fprintf(os.Stderr, "Error:\n-sec flag required for signing commit\n")
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

	ciphertext, err := encryptWithPublicKey(pubKeyPath, msg)
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
		rep, err := git.PlainOpen("./")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "\nError while reading:\n%v\n", err)
				os.Exit(1)
			}
			fmt.Println()

			openpgpEntity, err = getSigningEntity(secKeyPath, pass)
		} else {
			openpgpEntity = nil
		}

		_, err = CommitChanges(rep, []string{filename}, commit_message, openpgpEntity)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
		}
	}
}

func handleRemove(name string, commit_message string, noCommit bool, sign bool, secKeyPath string) {
	filename := name + ".gpg"
	if err := os.Remove(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing password:\n%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Password '%s' removed\n", name)

	if !noCommit {
		// Open repo
		rep, err := git.PlainOpen("./")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Git error:\n%v\n", err)
			os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "\nError while reading:\n%v\n", err)
				os.Exit(1)
			}
			fmt.Println()

			openpgpEntity, err = getSigningEntity(secKeyPath, pass)
		} else {
			openpgpEntity = nil
		}

		_, err = CommitChanges(rep, []string{filename}, commit_message, openpgpEntity)
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
