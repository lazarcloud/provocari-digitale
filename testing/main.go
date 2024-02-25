package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	_ "github.com/mattn/go-sqlite3"
)

func getBase64FileContent(filename string) string {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", filename, err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(fileContent)
}

func main() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	imageName := "cpp-executor:latest"

	envVars := map[string]string{
		"IS_STANDARD_IO":       "true",
		"TESTING_MODE":         "individualFiles",
		"SOURCE_BASE64":        getBase64FileContent("./data/main.cpp"),
		"NUMBER_OF_TEST_CASES": "3",
		"INPUT_0_BASE64":       getBase64FileContent("./input/1.in"),
		"INPUT_1_BASE64":       getBase64FileContent("./input/2.in"),
		"INPUT_2_BASE64":       getBase64FileContent("./input/3.in"),
		"OUTPUT_0_BASE64":      getBase64FileContent("./output/1.out"),
		"OUTPUT_1_BASE64":      getBase64FileContent("./output/2.out"),
		"OUTPUT_2_BASE64":      getBase64FileContent("./output/3.out"),
		// "SOLVE_ID":          database.GenerateUUID(),
		// "PROBLEM_ID":        "1234",
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
