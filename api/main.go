package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/lazarcloud/provocari-digitale/api/auth"
	"github.com/lazarcloud/provocari-digitale/api/auth/jwt"
	"github.com/lazarcloud/provocari-digitale/api/database"
	"github.com/lazarcloud/provocari-digitale/api/globals"
	"github.com/lazarcloud/provocari-digitale/api/utils"
	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	ID         string `json:"id"`
	Problem_ID string `json:"problem_id"`
	Code       string `json:"code"`
}

func main() {

	os.Remove("./database.sqlite")
	database.Connect()
	database.Populate()

	publicKey, err := jwt.CreateJWTWithClaims(globals.AuthAccessType, time.Hour*10000, "", globals.AuthRolePublic)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Public key: " + publicKey)

	r := mux.NewRouter()
	problemsRouter := r.PathPrefix("/api/problems").Subrouter()
	database.PrepareProblemsRouter(problemsRouter)
	problemsRouter.Use(auth.JWTMiddleware)

	authRouter := r.PathPrefix("/api/auth").Subrouter()
	database.PrepareAuthRouter(authRouter)
	authRouter.Use(auth.JWTMiddleware)

	solveRouterPublic := r.PathPrefix("/api/solve").Subrouter()
	solveRouterPublic.HandleFunc("/update/{id}", database.UpdateTestResultHandler).Methods("POST")
	solveRouterPublic.HandleFunc("/calculate/{id}/{status}", database.CalculateTestGroupFinalScoreHandler).Methods("GET")

	solveRouter := r.PathPrefix("/api/solve").Subrouter()
	solveRouter.HandleFunc("/submit/{id}", database.SolveHandler).Methods("POST")
	solveRouter.HandleFunc("/progress/{id}", database.GetSolveProgressHandler).Methods("GET") // TO DO: security for writing data and timestamps for test with status of created, to run, ran, finished and final score in test group
	solveRouter.HandleFunc("/{id}", database.GetUserSolvesHandler).Methods("GET")
	solveRouter.HandleFunc("/max_score/{id}", database.GetProblemMaxScoreHandler).Methods("GET")

	solveRouter.Use(auth.JWTMiddleware)

	fmt.Printf("Server is running on port %d...\n", globals.ApiPort)
	http.Handle("/", utils.CORSHandler.Handler(r))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", globals.ApiPort), nil))

}
