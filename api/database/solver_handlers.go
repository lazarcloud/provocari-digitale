package database

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/auth"
	"github.com/lazarcloud/provocari-digitale/api/utils"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userId := auth.GetUserId(r)

	if userId == "" {
		utils.RespondWithUnauthorized(w, "Unauthorized, please login before uploading code")
		return
	}
	// get code string from request body

	type RequestBody struct {
		Code string `json:"code"`
	}

	var body RequestBody
	err := utils.DecodeJSONBody(w, r, &body)
	if err != nil {
		utils.RespondWithError(w, "Failed to decode request body")
		fmt.Println(err.Error())
		return
	}

	codeBase64 := body.Code

	code, err := base64.StdEncoding.DecodeString(codeBase64)

	if err != nil {
		utils.RespondWithError(w, "Failed to decode base64 code")
		fmt.Println(err.Error())
		return
	}

	ok, testGroupId, err := CreateTestContainer(id, string(code), userId)

	if err != nil {
		utils.RespondWithError(w, "Failed to create test container")
		fmt.Println(err.Error())
		return
	}

	if !ok {
		utils.RespondWithError(w, "Failed to create test container")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status":        "ok",
		"problem_id":    id,
		"test_group_id": testGroupId,
	})
}

func GetUserSolvesHandler(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserId(r)

	if userId == "" {
		utils.RespondWithUnauthorized(w, "Unauthorized, please login before uploading code")
		return
	}

	solves, err := GetSolvesByUserId(userId)

	if err != nil {
		utils.RespondWithError(w, "Failed to retrieve solves")
		fmt.Println(err.Error())
		return
	}

	// if solve is null return empty array
	if solves == nil {
		utils.RespondWithSuccess(w, map[string]interface{}{
			"status": "ok",
			"solves": []Solve{},
		})
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status": "ok",
		"solves": solves,
	})
}

type Solve struct {
	ID         string `json:"id"`
	ProblemID  string `json:"problem_id"`
	FinalScore string `json:"final_score"`
	MaxScore   string `json:"max_score"`
	TestCount  int    `json:"test_count"`
	CreatedAt  string `json:"created_at"`
}

func GetSolvesByUserId(userId string) ([]Solve, error) {
	rows, err := DB.Query("SELECT id, problem_id, final_score, max_score, test_count, created_at FROM tests_groups WHERE user_id = ? ORDER BY created_at DESC", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var solves []Solve

	for rows.Next() {
		var solve Solve
		err := rows.Scan(&solve.ID, &solve.ProblemID, &solve.FinalScore, &solve.MaxScore, &solve.TestCount, &solve.CreatedAt)
		if err != nil {
			return nil, err
		}
		solves = append(solves, solve)
	}

	return solves, nil
}
