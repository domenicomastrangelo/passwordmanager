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
	)

	fmt.Println()
	fmt.Println("Adding password for: " + element)
	fmt.Println()

	fmt.Print("Please enter your element: ")

	if password, err = terminal.ReadPassword(0); err != nil {
		log.Fatalln("Could not read password from stdin")
	}

	if encryptedPassword, err = encrypt([]byte(userPasswordClear), password); err != nil {
		fmt.Println()
		fmt.Println(err.Error())
	}

	addElement("password", string(element), string(base64Encode(encryptedPassword)))
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

func get(elementType string, value string) bool {
	switch elementType {
	case "password":
		getPassword(value)
	}
	return false
}

func getPassword(element string) {
	var (
		password []byte
		err      error
	)

	fmt.Println()
	fmt.Println("Getting element for: " + element)

	rowElements := getElements("password", element)

	for _, rowElement := range rowElements {
		base64EncodedPassword := []byte(rowElement.Value)

		if password, err = base64Decode(string(base64EncodedPassword)); err != nil {
			fmt.Println()
			log.Fatalln("Could not decode base64 encoded password")
		}

		if password, err = decrypt([]byte(userPasswordClear), password); err != nil {
			fmt.Println()
			log.Fatalln("Could not decode decrypt password")
		}

		fmt.Println()
		fmt.Println("Element is: " + string(password))
	}
}

func base64Encode(key []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(key))
}

func base64Decode(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(key)
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
