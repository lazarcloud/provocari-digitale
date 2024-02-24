package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
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
