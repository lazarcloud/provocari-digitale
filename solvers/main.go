package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Output struct {
	MaxMemory string `json:"max_memory"`
	TimeTaken string `json:"time_taken"`
	Output    string `json:"output"`
	Error     string `json:"error"`
	ProblemID string `json:"problem_id"`
}

func finish(answer Output) {
	answer.ProblemID = os.Getenv("PROBLEM_ID")
	solveID := os.Getenv("SOLVE_ID")
	api := "http://host.docker.internal:8080/api/submit/" + solveID

	jsonData, err := json.Marshal(answer)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(api, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected response status:", resp.Status)
		return
	}

	fmt.Println(answer)
}

func main() {
	// Decode the base64-encoded executable
	executableBase64 := os.Getenv("EXECUTABLE_BASE64")
	executable, err := base64.StdEncoding.DecodeString(executableBase64)
	if err != nil {
		finish(Output{Error: "Error decoding base64"})
		return
	}

	// Write the decoded executable to a file
	err = os.WriteFile("executable", executable, 0755) // Set executable permission
	if err != nil {
		finish(Output{Error: "Error writing executable"})
		return
	}

	// Run the executable and capture stdout and stderr
	cmd := exec.Command("./executable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		finish(Output{Error: "Error running executable"})
		return
	}

	// Parse output to extract max memory and time taken
	timeLines := strings.Split(string(output), "\n")
	var maxMemory, timeTaken string
	for _, line := range timeLines {
		if strings.HasPrefix(line, "Max Memory:") {
			maxMemory = strings.TrimSpace(strings.TrimPrefix(line, "Max Memory:"))
		} else if strings.HasPrefix(line, "Elapsed Time:") {
			timeTaken = strings.TrimSpace(strings.TrimPrefix(line, "Elapsed Time:"))
		}
	}

	finish(Output{
		MaxMemory: maxMemory,
		TimeTaken: timeTaken,
		Output:    string(output),
		Error:     "",
	})

	// Clean up temporary file
	os.Remove("executable")
}
