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
	// Decode the base64-encoded C++ source code
	cppSourceBase64 := os.Getenv("CPP_SOURCE_BASE64")
	cppSource, err := base64.StdEncoding.DecodeString(cppSourceBase64)
	if err != nil {
		finish(Output{Error: "Error decoding base64"})
		return
	}

	// Write the decoded C++ source code to a file
	err = os.WriteFile("source.cpp", cppSource, 0644)
	if err != nil {
		finish(Output{Error: "Error writing source.cpp"})
		return
	}

	// Compile the C++ source code
	cmd := exec.Command("g++", "source.cpp", "-o", "output")
	output, err := cmd.CombinedOutput()
	if err != nil {
		finish(Output{Error: "Compilation failed: " + string(output)})
		return
	}

	// Run the compiled program and capture stdout and stderr
	cmd = exec.Command("./output")
	output, err = cmd.CombinedOutput()
	if err != nil {
		finish(Output{Error: "Error running program: " + string(output)})
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

	// Clean up temporary files
	os.Remove("source.cpp")
	os.Remove("output")
}
