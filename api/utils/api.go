package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, "Internal Server Error: "+err.Error())
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		RespondWithError(w, "Internal Server Error: "+err.Error())
		return
	}
}

func RespondWithError(w http.ResponseWriter, message string) {
	RespondWithJson(w, http.StatusInternalServerError, map[string]interface{}{"error": message})
}

func RespondWithSuccess(w http.ResponseWriter, message map[string]interface{}) {
	RespondWithJson(w, http.StatusOK, message)
}

func RespondWithUnauthorized(w http.ResponseWriter, message string) {
	RespondWithJson(w, http.StatusUnauthorized, map[string]interface{}{"error": message})
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
