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
