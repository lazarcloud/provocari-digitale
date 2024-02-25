package database

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/utils"
)

func CreateTestResult(test_id string, test_group_id string) (string, error) {
	id := GenerateUUID()
	_, err := DB.Exec("INSERT INTO tests_results (id, test_id, test_group_id, max_memory, time_taken, output, error) VALUES (?, ?, ?, ?, ?, ?, ?)", id, test_id, test_group_id, "", "", "", "", "")
	if err != nil {
		return "", err
	}
	return id, nil
}
func CreateTestGroup(problem_id string, user_id string) (string, error) {
	id := GenerateUUID()
	_, err := DB.Exec("INSERT INTO tests_groups (id, problem_id, user_id) VALUES (?, ?, ?)", id, problem_id, user_id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func UpdateTestResult(id string, max_memory string, time_taken string, output string, error string, correct bool) error {
	_, err := DB.Exec("UPDATE tests_results SET max_memory=?, time_taken=?, output=?, error=?, correct=? WHERE id=?", max_memory, time_taken, output, error, correct, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTestResultHandler(w http.ResponseWriter, r *http.Request) { //TO DO: access token
	vars := mux.Vars(r)
	testId := vars["id"]

	type RequestBody struct {
		MaxMemory string `json:"max_memory"`
		TimeTaken string `json:"time_taken"`
		Output    string `json:"output"`
		Error     string `json:"error"`
		Correct   bool   `json:"correct"`
	}

	var body RequestBody
	err := utils.DecodeJSONBody(w, r, &body)
	if err != nil {
		utils.RespondWithError(w, "Failed to decode request body")
		fmt.Println(err.Error())
		return
	}

	err = UpdateTestResult(testId, body.MaxMemory, body.TimeTaken, body.Output, body.Error, body.Correct)

	if err != nil {
		utils.RespondWithError(w, "Failed to update test result")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status":  "ok",
		"test_id": testId,
	})
}

func GetSolveProgressHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testId := vars["id"]

	type TestResult struct {
		MaxMemory string `json:"max_memory"`
		TimeTaken string `json:"time_taken"`
		Output    string `json:"output"`
		Error     string `json:"error"`
		Correct   bool   `json:"correct"`
	}

	var testResults []TestResult

	rows, err := DB.Query("SELECT max_memory, time_taken, output, error, correct FROM tests_results WHERE test_group_id=?", testId)
	if err != nil {
		utils.RespondWithError(w, "Failed to query database")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var testResult TestResult
		err = rows.Scan(&testResult.MaxMemory, &testResult.TimeTaken, &testResult.Output, &testResult.Error, &testResult.Correct)
		if err != nil {
			utils.RespondWithError(w, "Failed to scan row")
			return
		}
		testResults = append(testResults, testResult)
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status":  "ok",
		"results": testResults,
	})

}
