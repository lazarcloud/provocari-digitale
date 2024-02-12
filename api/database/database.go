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
		_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS tests (
			id BLOB PRIMARY KEY NOT NULL,
			problem_id BLOB NOT NULL,
			max_memory TEXT NOT NULL,
			time_taken TEXT NOT NULL,
			output TEXT NOT NULL,
			error TEXT NOT NULL
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

func InsertSolve(id string, problemID string, maxMemory string, timeTaken string, output string, error string) {
	_, err := DB.Exec("INSERT INTO tests (id, problem_id, max_memory, time_taken, output, error) VALUES (?, ?, ?, ?, ?, ?)", id, problemID, maxMemory, timeTaken, output, error)
	if err != nil {
		log.Fatal(err)
	}
}
