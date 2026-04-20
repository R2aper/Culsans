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
func encryptWithPublicKey(pubKeyPath string, msg []byte) ([]byte, error) {
	pubKeyArmored, err := readKeyFromFile(pubKeyPath)
	if err != nil {
		return nil, err
	}

	// Create pub key
	publicKey, err := crypto.NewKeyFromArmored(pubKeyArmored)
	if err != nil {
		return nil, err
	}

	// Set up PGP profile
	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().
		Recipient(publicKey).
		New()
	if err != nil {
		return nil, err
	}

	// Encryption
	pgpMessage, err := encHandle.Encrypt(msg)
	if err != nil {
		return nil, err
	}

	return pgpMessage.ArmorBytes()
}

func decryptWithPrivateKey(privKeyPath string, passphrase []byte, armoredMsg []byte) ([]byte, error) {
	privKeyArmored, err := readKeyFromFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.NewKeyFromArmored(privKeyArmored)
	defer privateKey.ClearPrivateParams() // !!Clearing private stuff
	if err != nil {
		return nil, err
	}

	// Unlocking
	unlockedKey, err := privateKey.Unlock(passphrase)
	if err != nil {
		return nil, fmt.Errorf("Incorrect password or damaged key: %w", err)
	}

	// Set up PGP profile
	pgp := crypto.PGP()
	decHandle, err := pgp.Decryption().
		DecryptionKey(unlockedKey).
		New()
	if err != nil {
		return nil, err
	}

	// Decryption
	decrypted, err := decHandle.Decrypt(armoredMsg, crypto.Armor)
	if err != nil {
		return nil, err
	}

	return decrypted.Bytes(), nil
}
