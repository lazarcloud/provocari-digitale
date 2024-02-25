package database

func GetProblemTestsCount(problemID string) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM tests WHERE problem_id = ?", problemID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
