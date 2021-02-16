package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	userPasswordClear string
	dbPath            *string
	action            *string
	elementType       *string
	element           *string
	actions           actionsType
	elements          elementTypes
)

func init() {
	actions = make(map[string]int)
	actions["add"] = 0
	actions["remove"] = 1

	elements = make(map[string]int)
	elements["password"] = 0
}

func main() {
	checkArgs()

	var (
		username []byte
		password []byte
		err      error
	)

	fmt.Print("Please enter your username: ")

	if _, err = fmt.Scanln(&username); err != nil {
		log.Fatalln("Could not read username from stdin")
	}

	fmt.Print("Please enter your password: ")

	if password, err = terminal.ReadPassword(0); err != nil {
		log.Fatalln("Could not read password from stdin")
	}

	fmt.Println()

	if checkPassword(password) {
		execute(*action)
	}

	fmt.Println("Done!")
}

func checkArgs() {
	dbPath = flag.String("path", "./db.sqlite", "--path=./db.sqlite")
	action = flag.String("action", "add", "-action=add")
	elementType = flag.String("elementType", "password", "-elementType=password")
	element = flag.String("element", "", "-element=https://www.google.com")

	flag.Parse()

	if _, ok := actions[*action]; !ok {
		log.Fatalln("Could find chosen action")
	}

	if _, ok := elements[*elementType]; !ok {
		log.Fatalln("Could find chosen element type")
	}

	if len(*element) <= 0 {
		log.Fatalln("Missing element")
	}
}

func checkPassword(password []byte) bool {
	userPasswordClear = string(password)

	return true
}

func execute(action string) {
	switch action {
	case "add":
		add(*elementType, *element)
	case "remove":
		remove(*elementType, *element)
	}
}
