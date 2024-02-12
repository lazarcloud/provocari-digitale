package database

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func Connect() {
	if _, err := os.Stat("./database.sqlite"); os.IsNotExist(err) {
		DB, err = sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}
		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS solves (
			id BLOB PRIMARY KEY NOT NULL,
			problem_id BLOB NOT NULL,
			code TEXT NOT NULL
		)`)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		DB, err = sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func InsertSolve(problemID string, code string) {
	_, err := DB.Exec("INSERT INTO solves (id, problem_id, code) VALUES (?, ?, ?)", GenerateUUID(), problemID, code)
	if err != nil {
		log.Println(err)
	}
}
