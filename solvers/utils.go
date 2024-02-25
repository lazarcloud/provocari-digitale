package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		finishExecution(fmt.Errorf("No %s provided", key))
	}
	return value
}
func decodeBase64(base64String string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		finishExecution(fmt.Errorf("Error decoding base64"))
	}
	return decoded
}
func writeFile(filename string, data []byte) {
	err := os.WriteFile(filename, data, 0755)
	if err != nil {
		finishExecution(fmt.Errorf("Error writing file"))
	}
}
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		finishExecution(fmt.Errorf("Error parsing int"))
	}
	return i
}

func compareFiles(file1, file2 string) bool {
	fmt.Println("Comparing files", file1, file2)
	content1, err := os.ReadFile(file1)
	if err != nil {
		fmt.Println("Error reading file1:", err)
		return false
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		fmt.Println("Error reading file2:", err)
		return false
	}
	fmt.Println("Comparing", string(content1), string(content2))
	return string(content1) == string(content2)
}

func saveTestResult(testID string, result bool, memory int64, executionTime time.Duration) {
	api := "http://host.docker.internal:8080/api/solve/update/" + testID
	type RequestBody struct {
		MaxMemory string `json:"max_memory"`
		TimeTaken string `json:"time_taken"`
		Output    string `json:"output"`
		Error     string `json:"error"`
		Correct   bool   `json:"correct"`
	}
	requestBody := RequestBody{
		MaxMemory: strconv.FormatInt(memory, 10),
		TimeTaken: executionTime.String(),
		Output:    "",
		Error:     "",
		Correct:   result,
	}

	jsonData, err := json.Marshal(requestBody)
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

	fmt.Println(requestBody)
}
