package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
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

func getItems(w http.ResponseWriter, r *http.Request) {
	var items []Item
	rows, err := database.DB.Query("SELECT id, problem_id, code FROM solves")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Problem_ID, &item.Code)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	json.NewEncoder(w).Encode(items)
}

func solve(w http.ResponseWriter, r *http.Request) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	imageName := "cpp-executor:latest"

	// read main.cpp form filesystem and convert it to base64

	fileContent, err := os.ReadFile("./data/main.cpp")
	if err != nil {
		fmt.Printf("Error reading main.cpp: %v\n", err)
		return
	}

	// Encode the file content to base64
	cppBase64 := base64.StdEncoding.EncodeToString(fileContent)

	envVars := map[string]string{
		"CPP_SOURCE_BASE64": cppBase64,
		"SOLVE_ID":          database.GenerateUUID(),
		"PROBLEM_ID":        "1",
	}
	var envList []string
	for key, value := range envVars {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	// Define container configuration
	containerConfig := &container.Config{
		Image: imageName,
		Env:   envList,
	}

	hostConfig := &container.HostConfig{
		// AutoRemove:  true,
		NetworkMode: "host",
		Resources:   container.Resources{
			// Memory: 64 * 1024 * 1024, // 256MB
			// CPUPeriod: 100000,1
			// CPUQuota:  10000, // 10ms (10% of a single CPU core)
		},
	}

	// Create the Docker container
	resp, err := dockerClient.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		fmt.Printf("Error creating Docker container: %v\n", err)
		return
	}
	// measure its stats

	fmt.Println("Docker container created successfully!")
	fmt.Printf("Container ID: %s\n", resp.ID)

	// Start the Docker container
	err = dockerClient.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		fmt.Printf("Error starting Docker container: %v\n", err)
		return
	}

	fmt.Println("Docker container started successfully!")
}

func main() {

	publicKey, err := jwt.CreateJWTWithClaims(globals.AuthAccessType, time.Hour*10000, "", globals.AuthRolePublic)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Public key: " + publicKey)

	os.Remove("./database.sqlite")
	database.Connect()
	database.Populate()
	r := mux.NewRouter()
	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/solve", solve).Methods("GET")

	r.HandleFunc("/api/submit/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		fmt.Println(id)

		//get json body
		type Output struct {
			MaxMemory  string `json:"max_memory"`
			TimeTaken  string `json:"time_taken"`
			Output     string `json:"output"`
			Error      string `json:"error"`
			Problem_ID string `json:"problem_id"`
		}

		var output Output
		err := json.NewDecoder(r.Body).Decode(&output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		database.InsertSolve(id, output.Problem_ID, output.MaxMemory, output.TimeTaken, output.Output, output.Error)

		fmt.Println(output)
		w.Write([]byte("OK"))
	}).Methods("POST")

	problemsRouter := r.PathPrefix("/api/problems").Subrouter()
	database.PrepareProblemsRouter(problemsRouter)
	problemsRouter.Use(auth.JWTMiddleware)

	authRouter := r.PathPrefix("/api/auth").Subrouter()
	database.PrepareAuthRouter(authRouter)
	authRouter.Use(auth.JWTMiddleware)

	fmt.Printf("Server is running on port %d...\n", globals.ApiPort)
	http.Handle("/", utils.CORSHandler.Handler(r))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", globals.ApiPort), nil))

}
