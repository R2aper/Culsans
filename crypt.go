package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

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

// Getting entity from private key for signing commit
func getSigningEntity(secKeyPath string, passphrase []byte) (*openpgp.Entity, error) {
	// I don't know why v2.Entity is not compatible with v1.Entity, like wtf??
	secKeyArmored, err := readKeyFromFile(secKeyPath)
	if err != nil {
		return nil, err
	}

	entityList, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(secKeyArmored))
	if len(entityList) == 0 {
		return nil, fmt.Errorf("Couldn't get key from %s: %v\n", secKeyPath, err)
	}

	entity := entityList[0]

	err = entity.PrivateKey.Decrypt(passphrase)
	if err != nil {
		return nil, err
	}

	return entity, nil
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
