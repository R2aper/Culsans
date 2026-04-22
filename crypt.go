package main

import (
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

func readKeyFromFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getPublicKey(pubKeyPath string) (*crypto.Key, error) {
	pubKeyArmored, err := readKeyFromFile(pubKeyPath)
	if err != nil {
		return nil, err
	}

	return crypto.NewKey(pubKeyArmored)
}

func getUnlockedPrivateKey(secKeyPath string, passphrase []byte) (*crypto.Key, error) {
	privKeyArmored, err := readKeyFromFile(secKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.NewKey(privKeyArmored)
	defer privateKey.ClearPrivateParams()
	if err != nil {
		return nil, err
	}

	return privateKey.Unlock(passphrase)
}

// Encrypt message and return armord string
func encryptWithPublicKey(pubKeyPath string, msg []byte) ([]byte, error) {
	publicKey, err := getPublicKey(pubKeyPath)
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

	return pgpMessage.Bytes(), nil
}

func decryptWithPrivateKey(secKeyPath string, passphrase []byte, encryptedData []byte) ([]byte, error) {
	unlockedKey, err := getUnlockedPrivateKey(secKeyPath, passphrase)
	if err != nil {
		return nil, err
	}

	pgpMsg := crypto.NewPGPMessage(encryptedData)

	// Set up PGP profile
	pgp := crypto.PGP()
	decHandle, err := pgp.Decryption().
		DecryptionKey(unlockedKey).
		New()
	if err != nil {
		return nil, err
	}

	// Decryption
	decrypted, err := decHandle.Decrypt(pgpMsg.Bytes(), crypto.Auto)
	if err != nil {
		return nil, err
	}

	return decrypted.Bytes(), nil
}
