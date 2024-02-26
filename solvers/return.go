package main

import (
	"fmt"
	"os"
)

func finishExecution(error error) {
	// Send the execution time to the parent process
	// The parent process will read the time and print it

	// send the error to the parent process

	if error != nil {
		numberOfTestCases := parseInt(getEnv("NUMBER_OF_TEST_CASES"))
		for i := 0; i < numberOfTestCases; i++ {
			test_id := getEnv(fmt.Sprintf("TEST_%d_ID", i))
			saveTestResult(test_id, false, 0, 0, "error")
		}

		test_group_id := getEnv("TEST_GROUP_ID")
		calculateScores(test_group_id, error.Error())

		fmt.Println("Error:", error)
		os.Exit(1)
	}
}
