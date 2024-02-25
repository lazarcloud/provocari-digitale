package database

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/auth"
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

func CreateTestGroup(problem_id string, user_id string, maxScore string, testCount int) (string, error) {
	id := GenerateUUID()
	_, err := DB.Exec("INSERT INTO tests_groups (id, problem_id, user_id, max_score, test_count) VALUES (?, ?, ?, ?, ?)", id, problem_id, user_id, maxScore, testCount)
	if err != nil {
		return "", err
	}
	return id, nil
}
func UpdateTestResult(id string, max_memory string, time_taken string, output string, error string, correct bool) error {
	_, err := DB.Exec("UPDATE tests_results SET max_memory=?, time_taken=?, output=?, error=?, correct=?, status=? WHERE id=?", max_memory, time_taken, output, error, correct, "finished", id)
	if err != nil {
		return err
	}
	return nil
}

func CalculateTestGroupFinalScore(test_group_id string) error { //TO DO: optimize
	// go through all test results and tests score field and calculate final score
	var finalScore int
	var maxScore int
	rows, err := DB.Query("SELECT test_id, correct FROM tests_results WHERE test_group_id=?", test_group_id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var testId string
		var correct bool
		err = rows.Scan(&testId, &correct)
		if err != nil {
			return err
		}
		// get score from tests of that id
		var score int
		err = DB.QueryRow("SELECT score FROM tests WHERE id=?", testId).Scan(&score)
		if err != nil {
			return err
		}
		if correct {
			finalScore += score
		}
		maxScore += score

	}

	_, err = DB.Exec("UPDATE tests_groups SET final_score=?, max_score=? WHERE id=?", finalScore, maxScore, test_group_id)
	if err != nil {
		return err
	}
	return nil

}
func CalculateTestGroupFinalScoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	testGroupId := vars["id"]

	err := CalculateTestGroupFinalScore(testGroupId)

	if err != nil {
		utils.RespondWithError(w, "Failed to calculate final score")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status":        "ok",
		"test_group_id": testGroupId,
	})
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

	// check if user is owner of the test group
	type TestGroup struct {
		MaxScore   string `json:"max_score"`
		FinalScore string `json:"final_score"`
		ProblemID  string `json:"problem_id"`
		TestCount  int    `json:"test_count"`
		UserID     string `json:"user_id"`
	}
	var testGroup TestGroup

	err := DB.QueryRow("SELECT max_score, final_score, problem_id, test_count, user_id FROM tests_groups WHERE id=?", test_group_id).Scan(&testGroup.MaxScore, &testGroup.FinalScore, &testGroup.ProblemID, &testGroup.TestCount, &testGroup.UserID)
	if err != nil {
		utils.RespondWithError(w, "Failed to query database")
		return
	}

	currentUserId := auth.GetUserId(r)

	if currentUserId != testGroup.UserID {
		utils.RespondWithError(w, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	test_group_id := vars["id"]

	type TestResult struct {
		MaxMemory string `json:"max_memory"`
		TimeTaken string `json:"time_taken"`
		Output    string `json:"output"`
		Error     string `json:"error"`
		Correct   bool   `json:"correct"`
	}

	var testResults []TestResult

	rows, err := DB.Query("SELECT max_memory, time_taken, output, error, correct FROM tests_results WHERE test_group_id=?", test_group_id)
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
		"group":   testGroup,
	})

}
