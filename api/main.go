// package main

// import (
// 	"context"
// 	"encoding/base64"
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/client"
// )

// func main() {
// 	// Initialize the Docker client
// 	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		fmt.Printf("Error creating Docker client: %v\n", err)
// 		return
// 	}

// 	// Define the Docker image name
// 	imageName := "cpp-executor:latest"

// 	// read main.cpp form filesystem and convert it to base64

// 	fileContent, err := ioutil.ReadFile("main.cpp")
// 	if err != nil {
// 		fmt.Printf("Error reading main.cpp: %v\n", err)
// 		return
// 	}

// 	// Encode the file content to base64
// 	cppBase64 := base64.StdEncoding.EncodeToString(fileContent)

// 	// Encode the C++ source file to base64
// 	cppBase64 := os.Getenv("CPP_SOURCE_BASE64")
// 	if cppBase64 == "" {
// 		fmt.Println("CPP_SOURCE_BASE64 environment variable is not set")
// 		return
// 	}

// 	// Define environment variables to pass to the Docker container
// 	envVars := map[string]string{
// 		"CPP_SOURCE_BASE64": cppBase64,
// 	}

// 	// Prepare environment variables
// 	var envList []string
// 	for key, value := range envVars {
// 		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
// 	}

// 	// Define container configuration
// 	containerConfig := &types.ContainerCreateConfig{
// 		Image: imageName,
// 		Env:   envList,
// 	}

// 	// Create the Docker container
// 	container, err := dockerClient.ContainerCreate(context.Background(), containerConfig, nil, nil, "")
// 	if err != nil {
// 		fmt.Printf("Error creating Docker container: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Docker container created successfully!")
// 	fmt.Printf("Container ID: %s\n", container.ID)

// 	// Start the Docker container
// 	err = dockerClient.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		fmt.Printf("Error starting Docker container: %v\n", err)
// 		return
// 	}

//		fmt.Println("Docker container started successfully!")
//	}
package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
	}
}
