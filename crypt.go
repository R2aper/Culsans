package main

import (
	"fmt"
	"os"
	"path"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

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

func readKeyFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Encrypt message and return armord string
func encryptWithPublicKey(pubKeyPath string, msg string) (string, error) {
	pubKeyArmored, err := readKeyFromFile(pubKeyPath)
	if err != nil {
		return "", err
	}

	// Create pub key
	publicKey, err := crypto.NewKeyFromArmored(pubKeyArmored)
	if err != nil {
		return "", err
	}

	// Set up PGP profile
	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().
		Recipient(publicKey).
		New()
	if err != nil {
		return "", err
	}

	// Encryption
	pgpMessage, err := encHandle.Encrypt([]byte(msg))
	if err != nil {
		return "", err
	}

	return pgpMessage.Armor()
}

func decryptWithPrivateKey(privKeyPath string, passphrase string, armoredMsg string) (string, error) {
	privKeyArmored, err := readKeyFromFile(privKeyPath)
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.NewKeyFromArmored(privKeyArmored)
	defer privateKey.ClearPrivateParams() // !!Clearing private stuff
	if err != nil {
		return "", err
	}

	// Unlocking
	unlockedKey, err := privateKey.Unlock([]byte(passphrase))
	if err != nil {
		return "", fmt.Errorf("Incorrect password or damaged key: %w", err)
	}

	// Set up PGP profile
	pgp := crypto.PGP()
	decHandle, err := pgp.Decryption().
		DecryptionKey(unlockedKey).
		New()
	if err != nil {
		return "", err
	}

	// Decryption
	decrypted, err := decHandle.Decrypt([]byte(armoredMsg), crypto.Armor)
	if err != nil {
		return "", err
	}

	return decrypted.String(), nil
}
