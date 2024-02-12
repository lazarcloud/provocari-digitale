package database

import (
	"database/sql"
	"log"
	"os"
)

var DB *sql.DB

func createTable(sql string) {
	_, err := DB.Exec(sql)
	if err != nil {
		log.Fatal(err)
	}
}

func Connect() {
	if _, err := os.Stat("./database.sqlite"); os.IsNotExist(err) {
		DB, err = sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}

		createTable(`CREATE TABLE IF NOT EXISTS tests (
			id BLOB PRIMARY KEY NOT NULL,
			problem_id BLOB NOT NULL,
			max_memory TEXT NOT NULL,
			time_taken TEXT NOT NULL,
			output TEXT NOT NULL,
			error TEXT NOT NULL
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS users (
			id BLOB PRIMARY KEY NOT NULL,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS problems (
			id BLOB PRIMARY KEY NOT NULL,
			owner_id BLOB NOT NULL,
			max_memory TEXT NOT NULL,
			max_time TEXT NOT NULL,
			description TEXT NOT NULL
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS solve_sources (
			id BLOB PRIMARY KEY NOT NULL,
			problem_id BLOB NOT NULL,
			file BLOB NOT NULL
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS solve_compiled_sources (
			id BLOB PRIMARY KEY NOT NULL,
			source_id BLOB NOT NULL,
			file BLOB NOT NULL
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS compilation_tasks (
			id BLOB PRIMARY KEY NOT NULL,
			source_id BLOB NOT NULL
		)`)

	} else {
		DB, err = sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Populate() {
	_, err := DB.Exec("INSERT INTO problems (id, max_memory, max_time, description) VALUES (?, ?, ?, ?)", "1", "256", "1", "Test problem")
	if err != nil {
		log.Fatal(err)
	}
}

func InsertSolve(id string, problemID string, maxMemory string, timeTaken string, output string, error string) {
	_, err := DB.Exec("INSERT INTO tests (id, problem_id, max_memory, time_taken, output, error) VALUES (?, ?, ?, ?, ?, ?)", id, problemID, maxMemory, timeTaken, output, error)
	if err != nil {
		log.Fatal(err)
	}
}
