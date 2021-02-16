package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh/terminal"
)

type actionsType map[string]int
type elementTypes map[string]int

func add(elementType string, value string) bool {
	switch elementType {
	case "password":
		addPassword(value)
	}
	return false
}

func addPassword(element string) {
	var (
		password          []byte
		err               error
		encryptedPassword []byte
		decryptedPassword []byte
	)

	fmt.Println()
	fmt.Println("Adding password for: " + element)
	fmt.Println()

	fmt.Print("Please enter your password: ")

	if password, err = terminal.ReadPassword(0); err != nil {
		log.Fatalln("Could not read password from stdin")
	}

	if encryptedPassword, err = encrypt([]byte(userPasswordClear), password); err != nil {
		fmt.Println()
		fmt.Println(err.Error())
	}

	if decryptedPassword, _ = decrypt([]byte(userPasswordClear), encryptedPassword); err != nil {
		fmt.Println()
		fmt.Println(err.Error())
	}

	fmt.Println()
	fmt.Println(string(base64Encode(encryptedPassword)))
	fmt.Println(string(decryptedPassword))
	fmt.Println()

	saveElement(element, password)
}

func remove(elementType string, value string) bool {
	switch elementType {
	case "password":
		removePassword(value)
	}

	return false
}

func removePassword(element string) {
	fmt.Println()
	fmt.Println("Removing password for: " + element)
	fmt.Println()
}

func base64Encode(key []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(key))
}

func sha256Hash(key []byte) []byte {
	h := sha256.New()
	h.Write(key)

	return h.Sum(nil)
}

func encrypt(key, data []byte) ([]byte, error) {
	key = sha256Hash(key)

	var (
		blockCipher, err = aes.NewCipher(key)
		gcm              cipher.AEAD
	)

	if err != nil {
		return nil, err
	}

	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

func decrypt(key, data []byte) ([]byte, error) {
	key = sha256Hash(key)

	var (
		blockCipher cipher.Block
		gcm         cipher.AEAD
		plaintext   []byte
		err         error
	)

	if blockCipher, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	if gcm, err = cipher.NewGCM(blockCipher); err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	if plaintext, err = gcm.Open(nil, nonce, ciphertext, nil); err != nil {
		return nil, err
	}

	return plaintext, nil
}

func saveElement(element string, value []byte) {

}
