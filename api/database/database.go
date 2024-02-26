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
		createTable(`CREATE TABLE IF NOT EXISTS users (
			id BLOB PRIMARY KEY NOT NULL,
			created_at INTEGER DEFAULT (CAST(strftime('%s', 'now') AS INT)),
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			verified BOOLEAN NOT NULL DEFAULT FALSE
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS problems (
			id BLOB PRIMARY KEY NOT NULL,
			uses_standard_io BOOLEAN NOT NULL DEFAULT TRUE,
			test_mode TEXT DEFAULT 'NULL',
			input_file_name TEXT DEFAULT 'NULL',
			output_file_name TEXT DEFAULT 'NULL',
			title TEXT NOT NULL,
			owner_id BLOB NOT NULL,
			max_memory TEXT NOT NULL,
			max_time TEXT NOT NULL,
			description TEXT NOT NULL,
			FOREIGN KEY(owner_id) REFERENCES users(id)
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS tests (
			id BLOB PRIMARY KEY NOT NULL,
			problem_id BLOB NOT NULL,
			input BLOB NOT NULL,
			output BLOB NOT NULL,
			score INTEGER DEFAULT 0,
			FOREIGN KEY(problem_id) REFERENCES problems(id)
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS tests_results (
			id BLOB PRIMARY KEY NOT NULL,
			test_id BLOB NOT NULL,
			test_group_id BLOB NOT NULL,
			max_memory TEXT NOT NULL,
			time_taken TEXT NOT NULL,
			output BLOB NOT NULL,
			error BLOB NOT NULL,
			status TEXT DEFAULT 'waiting',
			correct BOOLEAN NOT NULL DEFAULT FALSE,
			FOREIGN KEY(test_id) REFERENCES tests(id),
			FOREIGN KEY(test_group_id) REFERENCES tests_groups(id)
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS tests_groups (
			id BLOB PRIMARY KEY NOT NULL,
			created_at INTEGER DEFAULT (CAST(strftime('%s', 'now') AS INT)),
			user_id BLOB NOT NULL,
			problem_id BLOB NOT NULL,
			final_score TEXT DEFAULT 'NULL',
			max_score TEXT DEFAULT 'NULL',
			test_count INTEGER DEFAULT 0,
			status TEXT DEFAULT 'waiting',
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(problem_id) REFERENCES problems(id)
		)`)
		createTable(`CREATE TABLE IF NOT EXISTS spam (
			user_id BLOB PRIMARY KEY NOT NULL,
			last_spam INTEGER DEFAULT (CAST(strftime('%s', 'now') AS INT))
		)`)
		// createTable(`CREATE TABLE IF NOT EXISTS solve_sources (
		// 	id BLOB PRIMARY KEY NOT NULL,
		// 	problem_id BLOB NOT NULL,
		// 	file BLOB NOT NULL
		// )`)
		// createTable(`CREATE TABLE IF NOT EXISTS solve_compiled_sources (
		// 	id BLOB PRIMARY KEY NOT NULL,
		// 	source_id BLOB NOT NULL,
		// 	file BLOB NOT NULL
		// )`)
		// createTable(`CREATE TABLE IF NOT EXISTS compilation_tasks (
		// 	id BLOB PRIMARY KEY NOT NULL,
		// 	source_id BLOB NOT NULL
		// )`)

	} else {
		DB, err = sql.Open("sqlite3", "./database.sqlite")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Populate() {
	statement, err := DB.Prepare("INSERT INTO users (id, email, password, username) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec("laz", "system", "", GenerateRandomUsername())
	if err != nil {
		log.Fatal(err.Error())
	}
	statement, err = DB.Prepare("INSERT INTO users (id, email, password, username) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec("8f3f3855-ff2a-42fd-a595-d14fe683b488", "lazar@lazar.lol", "27568c7bfb1fe49ece7cefed431a638c14ab8b65", GenerateRandomUsername())
	if err != nil {
		log.Fatal(err.Error())
	}
	problems := []struct {
		id               string
		maxMemory        string
		maxTime          string
		description      string
		title            string
		uses_standard_io bool
		test_mode        string
	}{
		{"1", "256", "1", "Test problem 1", "1 pb", true, "NULL"},
		{"2", "512", "2", "Test problem 2", "2 pb", true, "NULL"},
		{"3", "1024", "3", "Test problem 3", "3 pb", true, "NULL"},
		{"1234", "1024", "3", "Citește două numere întregi din cin și afișează suma lor în cout.", "A + B", true, "individualFiles"},
		{"1235", "1024", "3", "Citește două numere întregi din cin și afișează suma lor în cout.", "A + B", true, "individualFiles"},
	}

	for _, problem := range problems {
		_, err := DB.Exec("INSERT INTO problems (id, title, owner_id, max_memory, max_time, description, uses_standard_io, test_mode) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", problem.id, problem.title, "laz", problem.maxMemory, problem.maxTime, problem.description, problem.uses_standard_io, problem.test_mode)
		if err != nil {
			log.Fatal(err)
		}
	}

	tests := []struct {
		problemID string
		input     string
		output    string
		score     int
	}{
		{"1234", "1 2", "3", 1},
		{"1234", "5 4", "9", 1},
		{"1234", "1 5", "6", 1},
		{"1234", "4 7", "11", 1},
		{"1234", "3 7", "10", 1},
		{"1234", "1 9", "10", 1},
		{"1234", "5 8", "13", 1},
		{"1234", "8 9", "17", 1},
		{"1234", "21 4", "25", 1},
		{"1234", "5 5", "10", 1},
		{"1235", "1 2", "3", 1},
		{"1235", "5 4", "9", 1},
	}
	for _, test := range tests {
		_, err := DB.Exec("INSERT INTO tests (id, problem_id, input, output, score) VALUES (?, ?, ?, ?, ?)", GenerateUUID(), test.problemID, test.input, test.output, test.score)
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
