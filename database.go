package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db  *sql.DB
	err error
)

type elementData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func init() {
	if db, err = sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/test"); err != nil {
		log.Fatalln("Could not connect to db")
	}

	provisionDatabase()
}

func provisionDatabase() {
	createUsersTables()
	createElementTypesTables()
	createElementTables()
}

func createElementTypesTables() {
	var (
		err   error
		query string
	)

	query = `
		CREATE TABLE IF NOT EXISTS element_types (
			id INT AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,

			UNIQUE(name),
			PRIMARY KEY(id)
		);
	`

	if _, err = db.Exec(query); err != nil {
		log.Fatalln(err.Error())
	}

	query = `
		INSERT IGNORE INTO element_types(name) values("password")
	`

	if _, err = db.Exec(query); err != nil {
		log.Fatalln(err.Error())
	}
}

func createElementTables() {
	var (
		err error
	)

	query := `
		CREATE TABLE IF NOT EXISTS elements (
			id INT PRIMARY KEY AUTO_INCREMENT,
			user_id INT NOT NULL,
			element_type INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			value LONGTEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

			FOREIGN KEY(element_type) REFERENCES element_types(id),
			FOREIGN KEY(user_id) REFERENCES users(id)
		)
	`

	if _, err = db.Exec(query); err != nil {
		log.Fatalln(err.Error())
	}
}

func createUsersTables() {
	var (
		err error
	)

	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`

	if _, err = db.Exec(query); err != nil {
		log.Fatalln(err.Error())
	}
}

func addElement(elementType string, elementName string, value string) {
	var (
		err           error
		elementTypeID string
		row           *sql.Row
		query         string
	)

	query = `
		SELECT id FROM element_types WHERE name = ?
	`

	row = db.QueryRow(query, elementType)

	if err = row.Scan(&elementTypeID); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}

	userID := getUserID()

	query = `
		INSERT INTO elements(element_type, user_id, name, value) values(?, ?, ?, ?)
	`

	if _, err = db.Exec(query, elementTypeID, userID, elementName, value); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}
}

func getElements(elementType string, elementName string) []elementData {
	var (
		err           error
		elementTypeID string
		row           *sql.Row
		rows          *sql.Rows
		query         string
		edSingle      elementData
		ed            []elementData
	)

	query = `
		SELECT id FROM element_types WHERE name = ?
	`

	row = db.QueryRow(query, elementType)

	if err = row.Scan(&elementTypeID); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}

	query = `
		SELECT name, value FROM elements where element_type = ? and name = ? and user_id = ?
	`

	if rows, err = db.Query(query, elementTypeID, elementName, getUserID()); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}

	for rows.Next() {

		if err = rows.Scan(&edSingle.Name, &edSingle.Value); err != nil {
			log.Println()
			log.Fatalln(err.Error())
		}

		ed = append(ed, edSingle)
	}

	return ed
}

func getUserID() string {
	var (
		err    error
		userID string
		row    *sql.Row
		query  string
	)

	query = `
		SELECT id FROM users WHERE name = ?
	`

	row = db.QueryRow(query, username)

	if err = row.Scan(&userID); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}

	return userID
}

func checkUser(username []byte, password []byte) bool {
	userPasswordClear = string(password)

	var (
		err               error
		pass              string
		decodedPassword   []byte
		encryptedPassword []byte
		row               *sql.Row
		query             string
	)

	query = `
		SELECT password FROM users WHERE name = ?
	`

	row = db.QueryRow(query, username)

	if err = row.Scan(&pass); err == sql.ErrNoRows {
		fmt.Println()
		log.Println("Could not find user")
		log.Println("Creating it")

		if encryptedPassword, err = encrypt(password, password); err != nil {
			log.Println()
			log.Fatalln(err.Error())
		}

		pass = string(base64Encode(encryptedPassword))

		createUser(username, pass)
	} else if err != nil {
		log.Fatalln(err.Error())
	}

	if decodedPassword, err = base64Decode(pass); err != nil {
		log.Fatalln(err.Error())
	}

	if _, err = decrypt(password, []byte(decodedPassword)); err != nil {
		fmt.Println()
		log.Println("Login failed")
		return false
	}

	return true
}

func createUser(username []byte, password string) {
	var (
		err   error
		query string
	)

	query = `
		INSERT INTO users(name, password) values(?, ?)
	`

	if _, err = db.Exec(query, username, password); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}
}
