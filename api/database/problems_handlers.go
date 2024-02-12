package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/auth"
	"github.com/lazarcloud/provocari-digitale/api/globals"
	"github.com/lazarcloud/provocari-digitale/api/utils"
	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateProblemHandler handles POST requests to create a new problem.
func CreateProblemHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserId(r)
	if userID == "" {
		writeJSONError(w, "You are not allowed to create a problem, please log in", http.StatusForbidden)
		return
	}
	var problem Problem
	err := json.NewDecoder(r.Body).Decode(&problem)
	if err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	problem.OwnerID = userID

	err = CreateProblem(problem)
	if err != nil {
		writeJSONError(w, "Failed to create problem", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Problem created successfully")
}

// GetProblemsHandler handles GET requests to retrieve all problems.
func GetProblemsHandler(w http.ResponseWriter, r *http.Request) {
	defaultPage := 1
	defaultPageSize := 10

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = defaultPageSize
	}

	problems, err := GetAllProblems(page, pageSize)
	if err != nil {
		writeJSONError(w, "Failed to retrieve problems", http.StatusInternalServerError)
		return
	}

	totalProblems, err := GetTotalProblems()
	if err != nil {
		writeJSONError(w, "Failed to retrieve total number of problems", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalProblems) / float64(pageSize)))

	response := struct {
		CurrentPage   int       `json:"current_page"`
		TotalPages    int       `json:"total_pages"`
		PageSize      int       `json:"page_size"`
		TotalProblems int       `json:"total_problems"`
		Problems      []Problem `json:"problems"`
	}{
		CurrentPage:   page,
		TotalPages:    totalPages,
		PageSize:      pageSize,
		TotalProblems: totalProblems,
		Problems:      problems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProblemHandler handles PUT requests to update an existing problem by ID.
func UpdateProblemHandler(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserId(r)
	vars := mux.Vars(r)
	problemID := vars["id"]

	canDel, _ := IsProblemOwner(userID, problemID)
	if !canDel {
		writeJSONError(w, "You are not allowed to delete this problem", http.StatusForbidden)
		return
	}

	var updatedProblem Problem
	err := json.NewDecoder(r.Body).Decode(&updatedProblem)
	if err != nil {
		writeJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = UpdateProblem(problemID, updatedProblem)
	if err != nil {
		writeJSONError(w, "Failed to update problem", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Problem updated successfully")
}

// DeleteProblemHandler handles DELETE requests to delete a problem by ID.
func DeleteProblemHandler(w http.ResponseWriter, r *http.Request) {
	role := auth.GetRole(r)
	userID := auth.GetUserId(r)
	vars := mux.Vars(r)
	problemID := vars["id"]

	if role != globals.AuthRoleService {
		canDel, _ := IsProblemOwner(userID, problemID)
		if !canDel {
			writeJSONError(w, "You are not allowed to delete this problem", http.StatusForbidden)
			return
		}
	}

	err := DeleteProblem(problemID)
	if err != nil {
		writeJSONError(w, "Failed to delete problem", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Problem deleted successfully")
}

// GetProblemHandler handles GET requests to retrieve a single problem by ID.
func GetProblemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	problem, err := GetProblemByID(problemID)
	if err != nil {
		if err == sql.ErrNoRows {
			writeEmptyJSON(w)
			return
		}
		writeJSONError(w, "Failed to retrieve problem", http.StatusInternalServerError)
		return
	}

	// Return problem as JSON object
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

// writeJSONError writes a JSON error response with the given message and status code.
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	utils.RespondWithError(w, message)
}

// writeEmptyJSON writes an empty JSON response.
func writeEmptyJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}

func PrepareProblemsRouter(problemsRouter *mux.Router) {
	problemsRouter.HandleFunc("", GetProblemsHandler).Methods("GET")
	problemsRouter.HandleFunc("", CreateProblemHandler).Methods("POST")
	problemsRouter.HandleFunc("/{id}", GetProblemHandler).Methods("GET")
	problemsRouter.HandleFunc("/{id}", UpdateProblemHandler).Methods("PUT")
	problemsRouter.HandleFunc("/{id}", DeleteProblemHandler).Methods("DELETE")
}
