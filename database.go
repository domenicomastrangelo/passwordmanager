package main

import (
	"database/sql"
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
			element_type INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			value LONGTEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

			FOREIGN KEY(element_type) REFERENCES element_types(id)
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

	query = `
		INSERT INTO elements(element_type, name, value) values(?, ?, ?)
	`

	if _, err = db.Exec(query, elementTypeID, elementName, value); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}
}

func getElements(elementType string, elementName string) elementData {
	var (
		err           error
		elementTypeID string
		row           *sql.Row
		query         string
		ed            elementData
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
		SELECT name, value FROM elements where element_type = ? and name = ?
	`

	row = db.QueryRow(query, elementTypeID, elementName)

	if err = row.Scan(&ed.Name, &ed.Value); err != nil {
		log.Println()
		log.Fatalln(err.Error())
	}

	return ed
}
