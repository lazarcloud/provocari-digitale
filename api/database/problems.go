package database

import (
	"database/sql"
	"errors"
)

type Problem struct {
	ID          string `json:"id"`
	OwnerID     string `json:"owner_id"`
	MaxMemory   string `json:"max_memory"`
	MaxTime     string `json:"max_time"`
	Description string `json:"description"`
}

// GetAllProblems returns a slice of all problems with pagination support.
func GetAllProblems(page, pageSize int) ([]Problem, error) {
	offset := (page - 1) * pageSize

	rows, err := DB.Query("SELECT id, owner_id, max_memory, max_time, description FROM problems LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []Problem
	for rows.Next() {
		var problem Problem
		err := rows.Scan(&problem.ID, &problem.OwnerID, &problem.MaxMemory, &problem.MaxTime, &problem.Description)
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
	_, err := DB.Exec("INSERT INTO problems (id, owner_id, max_memory, max_time, description) VALUES (?, ?, ?, ?, ?)",
		problem.ID, problem.OwnerID, problem.MaxMemory, problem.MaxTime, problem.Description)
	if err != nil {
		return err
	}
	return nil
}

// GetProblemByID retrieves a problem by its ID.
func GetProblemByID(id string) (*Problem, error) {
	var problem Problem
	err := DB.QueryRow("SELECT id, owner_id, max_memory, max_time, description FROM problems WHERE id = ?", id).
		Scan(&problem.ID, &problem.OwnerID, &problem.MaxMemory, &problem.MaxTime, &problem.Description)
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
	_, err := DB.Exec("UPDATE problems SET owner_id=?, max_memory=?, max_time=?, description=? WHERE id=?",
		updatedProblem.OwnerID, updatedProblem.MaxMemory, updatedProblem.MaxTime, updatedProblem.Description, id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProblem deletes a problem by its ID.
func DeleteProblem(id string) error {
	_, err := DB.Exec("DELETE FROM problems WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}
