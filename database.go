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
		err error
	)

	query := `
		CREATE TABLE IF NOT EXISTS element_types (
			id INT AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,

			UNIQUE(name),
			PRIMARY KEY(id)
		)
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
