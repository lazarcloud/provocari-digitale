package database

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

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

	//check spam table for last accessed

	type SpamType struct {
		user_id   string
		last_spam int64
	}

	var spam SpamType

	err := DB.QueryRow("SELECT user_id, last_spam FROM spam WHERE user_id = ?", userId).Scan(&spam.user_id, &spam.last_spam)
	var justInserted bool = false
	if err != nil {
		// if user is not in spam table, add him
		if err == sql.ErrNoRows {
			_, err := DB.Exec("INSERT INTO spam (user_id) VALUES (?)", userId)
			if err != nil {
				utils.RespondWithError(w, "Failed to insert user into spam table")
				fmt.Println(err.Error())
				return
			}
			spam.last_spam = time.Now().Unix()
			justInserted = true
		} else {
			utils.RespondWithError(w, "Failed to query spam table")
			fmt.Println(err.Error())
			return
		}
	}

	// update last spam in db

	_, err = DB.Exec("UPDATE spam SET last_spam = ? WHERE user_id = ?", time.Now().Unix(), userId)

	if err != nil {
		utils.RespondWithError(w, "Failed to update last spam")
		fmt.Println(err.Error())
		return
	}

	deltaTime := time.Now().Unix() - spam.last_spam

	fmt.Println(time.Now().Unix())
	fmt.Println(spam.last_spam)

	if deltaTime < 10 && !justInserted {
		utils.RespondWithError(w, "You are spamming the server, please wait 10 seconds before submitting another solve")
		return
	}

	// get code string from request body

	type RequestBody struct {
		Code string `json:"code"`
	}

	var body RequestBody
	err = utils.DecodeJSONBody(w, r, &body)
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
	vars := mux.Vars(r)
	problem_id := vars["id"]

	userId := auth.GetUserId(r)

	if userId == "" {
		utils.RespondWithUnauthorized(w, "Unauthorized, please login before uploading code")
		return
	}

	solves, err := GetSolvesByUserId(userId, problem_id)

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

func GetSolvesByUserId(userId string, problem_id string) ([]Solve, error) {
	rows, err := DB.Query("SELECT id, problem_id, final_score, max_score, test_count, created_at FROM tests_groups WHERE user_id = ? AND problem_id = ? ORDER BY created_at DESC", userId, problem_id)
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


createTable(`CREATE TABLE IF NOT EXISTS tests_groups (
	id BLOB PRIMARY KEY NOT NULL,
	created_at INTEGER DEFAULT (CAST(strftime('%s', 'now') AS INT)),
	user_id BLOB NOT NULL,
	problem_id BLOB NOT NULL,
	final_score TEXT DEFAULT 'NULL',
	max_score TEXT DEFAULT 'NULL',
	test_count INTEGER DEFAULT 0,
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY(problem_id) REFERENCES problems(id)
)`)

func GetMaxScorePerProblemPerUser(user_id string, problem_id string) int{
	
	rows, err := DB.Query("SELECT max_score FROM tests_groups WHERE user_id = ? AND problem_id = ?", user_id, problem_id)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var maxScores []string

	for rows.Next() {
		var maxScore string
		err := rows.Scan(&maxScore)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		maxScores = append(maxScores, maxScore)
	}

	maxScore := 0
	for _, score := range maxScores {
		// convert score to int
		if(score == "NULL") {
			continue
		}

		intScore := 0
		fmt.Sscanf(score, "%d", &intScore)

		if intScore > maxScore {
			maxScore = intScore
		}
	}
	return maxScore
}

func GetProblemMaxScoreHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	problem_id := vars["id"]

	userId := auth.GetUserId(r)

	if userId == "" {
		utils.RespondWithUnauthorized(w, "Unauthorized, please login before uploading code")
		return
	}

	maxScore := GetMaxScorePerProblemPerUser(userId, problem_id)

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status": "ok",
		"max_score": maxScore,
	})
}
