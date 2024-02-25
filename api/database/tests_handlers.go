package database

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/utils"
)

func GetProblemTestsCountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	count, err := GetProblemTestsCount(problemID)
	if err != nil {
		utils.RespondWithError(w, "Failed to get problem tests count")
		return
	}

	utils.RespondWithSuccess(w, map[string]interface{}{
		"status":     "ok",
		"problem_id": problemID,
		"count":      count,
	})
}
