package database

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/lazarcloud/provocari-digitale/api/globals"
)

func convertToBase64(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
func parseInt(s string) int64 {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0
	}
	return i
}
func CreateTestContainer(problemID string, solution string, usedId string) (bool, string, error) {

	fmt.Println("Creating Docker container...")

	fmt.Println("Problem ID:", problemID)
	fmt.Println("Solution:", solution)

	problemData, err := GetProblemByID(problemID)
	if err != nil {
		return false, "", err
	}

	fmt.Println("Problem data:", problemData)

	fmt.Println("Source code:", solution)

	// Create a new test container
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, "", err
	}

	environmentVariables := map[string]string{}

	environmentVariables["PROBLEM_ID"] = problemID
	if problemData.UsesStandardIO {
		environmentVariables["USES_STANDARD_IO"] = "true"
	} else {
		environmentVariables["USES_STANDARD_IO"] = "false"
	}
	environmentVariables["TEST_MODE"] = problemData.TestMode
	environmentVariables["SOURCE_BASE64"] = convertToBase64(solution)
	environmentVariables["PROBLEM_INPUT_FILE"] = problemData.InputFileName
	environmentVariables["PROBLEM_OUTPUT_FILE"] = problemData.OutputFileName
	testCount, err := GetProblemTestsCount(problemID) // TO DO: optimize this in one sql query
	if err != nil {
		return false, "", err
	}
	environmentVariables["NUMBER_OF_TEST_CASES"] = fmt.Sprintf("%d", testCount)

	// Get all tests for the problem
	tests, err := GetAllTests(problemID)
	if err != nil {
		return false, "", err
	}

	for i, test := range tests {
		environmentVariables[fmt.Sprintf("INPUT_%d_BASE64", i)] = convertToBase64(test.Input)
		environmentVariables[fmt.Sprintf("OUTPUT_%d_BASE64", i)] = convertToBase64(test.Output)
	}

	maxScore, err := GetProblemMaxScore(problemID)

	if err != nil {
		return false, "", err
	}

	testGroupId, err := CreateTestGroup(problemID, usedId, fmt.Sprintf("%d", maxScore), testCount, solution)
	if err != nil {
		return false, "", err
	}

	environmentVariables["TEST_GROUP_ID"] = testGroupId

	for i, test := range tests { //TO DO: login users
		test_id, err := CreateTestResult(test.ID, testGroupId)
		if err != nil {
			return false, "", err
		}
		environmentVariables[fmt.Sprintf("TEST_%d_ID", i)] = test_id
	}

	// max memory and max time
	environmentVariables["MAX_MEMORY"] = problemData.MaxMemory
	environmentVariables["MAX_TIME"] = problemData.MaxTime

	var envList []string
	for key, value := range environmentVariables {
		envList = append(envList, fmt.Sprintf("%s=%s", key, value))
	}

	// Define container configuration
	containerConfig := &container.Config{
		Image: globals.SolverImageName,
		Env:   envList,
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: "host",
		Resources: container.Resources{
			Memory: parseInt(problemData.MaxMemory) * 1024, //MB
			// CPUPeriod: 100000,1
			CPUQuota: 10000, // 10ms (10% of a single CPU core) //laptop 10000
		},
	}

	// Create the Docker container
	resp, err := dockerClient.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return false, "", err
	}
	// measure its stats

	fmt.Println("Docker container created successfully!")
	fmt.Printf("Container ID: %s\n", resp.ID)

	// Start the Docker container
	err = dockerClient.ContainerStart(context.Background(), resp.ID, container.StartOptions{})
	if err != nil {
		return false, "", err
	}

	fmt.Println("Docker container started successfully!")

	return true, testGroupId, nil
}
