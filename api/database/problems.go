package database

import (
	"database/sql"
	"errors"
)

type Problem struct {
	ID             string `json:"id"`
	OwnerID        string `json:"owner_id"`
	OwnerEmail     string `json:"owner_email"`
	Title          string `json:"title"`
	MaxMemory      string `json:"max_memory"`
	MaxTime        string `json:"max_time"`
	Description    string `json:"description"`
	UsesStandardIO bool   `json:"uses_standard_io"`
	TestMode       string `json:"test_mode"`
	InputFileName  string `json:"input_file_name"`
	OutputFileName string `json:"output_file_name"`
}

// GetAllProblems returns a slice of all problems with pagination support.
func GetAllProblems(page, pageSize int) ([]Problem, error) {
	offset := (page - 1) * pageSize

	rows, err := DB.Query(`
		SELECT p.id, p.title, p.owner_id, p.max_memory, p.max_time, p.description, p.uses_standard_io, p.test_mode, p.input_file_name, p.output_file_name, u.email 
		FROM problems p 
		INNER JOIN users u ON p.owner_id = u.id 
		LIMIT ? OFFSET ?`, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var problem Problem
		err := rows.Scan(&problem.ID, &problem.Title, &problem.OwnerID, &problem.MaxMemory, &problem.MaxTime, &problem.Description, &problem.UsesStandardIO, &problem.TestMode, &problem.InputFileName, &problem.OutputFileName, &problem.OwnerEmail)
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

func GetTotalProblems() (int, error) {
	var total int
	err := DB.QueryRow("SELECT COUNT(*) FROM problems").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// CreateProblem creates a new problem.
func CreateProblem(problem Problem) error {
	_, err := DB.Exec("INSERT INTO problems (id, title, owner_id, max_memory, max_time, description) VALUES (?, ?, ?, ?, ?, ?)",
		problem.ID, problem.Title, problem.OwnerID, problem.MaxMemory, problem.MaxTime, problem.Description)
	if err != nil {
		return err
	}
	return nil
}

// GetProblemByID retrieves a problem by its ID.
func GetProblemByID(id string) (*Problem, error) {
	var problem Problem
	err := DB.QueryRow("SELECT id, title, owner_id, max_memory, max_time, description, uses_standard_io, test_mode, input_file_name, output_file_name FROM problems WHERE id = ?", id).
		Scan(&problem.ID, &problem.Title, &problem.OwnerID, &problem.MaxMemory, &problem.MaxTime, &problem.Description, &problem.UsesStandardIO, &problem.TestMode, &problem.InputFileName, &problem.OutputFileName)
	if err == sql.ErrNoRows {
		return nil, errors.New("problem not found")
	}
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// UpdateProblem updates an existing problem.
func UpdateProblem(id string, updatedProblem Problem) error {
	_, err := DB.Exec("UPDATE problems SET owner_id=?, title=?, max_memory=?, max_time=?, description=? WHERE id=?",
		updatedProblem.OwnerID, updatedProblem.Title, updatedProblem.MaxMemory, updatedProblem.MaxTime, updatedProblem.Description, id)
	if err != nil {
		return err
	}
	return nil
}

func IsProblemOwner(id string, userID string) (bool, error) {
	var ownerID string
	err := DB.QueryRow("SELECT owner_id FROM problems WHERE id=?", id).Scan(&ownerID)
	if err != nil {
		return false, err
	}
	return ownerID == userID, nil

}

// DeleteProblem deletes a problem by its ID.
func DeleteProblem(id string) error {
	_, err := DB.Exec("DELETE FROM problems WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

type Test struct {
	ID        string `json:"id"`
	ProblemID string `json:"problem_id"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	Count     int    `json:"count"`
}

func GetAllTests(problemID string) ([]Test, error) {
	rows, err := DB.Query("SELECT id, input, output, count FROM tests WHERE problem_id=?", problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []Test
	for rows.Next() {
		var test Test
		err := rows.Scan(&test.ID, &test.Input, &test.Output, &test.Count)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func GetAllTestsJSON(problemID string) ([]string, error) {
	rows, err := DB.Query("SELECT input, output FROM tests WHERE problem_id=?", problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []string
	for rows.Next() {
		var input, output string
		err := rows.Scan(&input, &output)
		if err != nil {
			return nil, err
		}
		tests = append(tests, input, output)
	}

	return tests, nil
}
