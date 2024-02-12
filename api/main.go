package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
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
		AutoRemove: true,
		Resources:  container.Resources{
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
	err = dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Printf("Error starting Docker container: %v\n", err)
		return
	}

	fmt.Println("Docker container started successfully!")
}
