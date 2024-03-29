package main

import (
	"fmt"
	"time"
)

var validTestingModes = []string{"individualFiles"}

func paseTimeDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		finishExecution(fmt.Errorf("Error parsing time duration"))
	}
	return d

}

func main() {
	// Get problem details

	isStandardIO := getEnv("USES_STANDARD_IO")
	testingMode := getEnv("TEST_MODE")
	problemInputFile := getEnv("PROBLEM_INPUT_FILE")
	problemOutputFile := getEnv("PROBLEM_OUTPUT_FILE")
	numberOfTestCases := parseInt(getEnv("NUMBER_OF_TEST_CASES"))
	maxMemory := parseInt64(getEnv("MAX_MEMORY"))
	maxTime := paseTimeDuration(getEnv("MAX_TIME") + "ms")
	// print the problem details
	fmt.Println("Testing mode:", testingMode)
	fmt.Println("Problem input file:", problemInputFile)
	fmt.Println("Problem output file:", problemOutputFile)
	fmt.Println("Number of test cases:", numberOfTestCases)

	// Check if the testing mode is valid
	valid := false
	for _, mode := range validTestingModes {
		if testingMode == mode {
			valid = true
			break
		}
	}
	if !valid {
		finishExecution(fmt.Errorf("Invalid testing mode"))
	}

	// Get the source code and write it to a file
	cppSourceBase64 := getEnv("SOURCE_BASE64")
	cppSource := decodeBase64(cppSourceBase64)
	fmt.Println("Source code decoded successfully")
	fmt.Println(string(cppSource))
	writeFile("./main.cpp", cppSource)
	fmt.Println("Source code written successfully")
	test_group_id := getEnv("TEST_GROUP_ID")

	// Compile the source code
	calculateScores(test_group_id, "compiling")
	compileCPP("./main.cpp", "executable")
	fmt.Println("Source code compiled successfully")
	time.Sleep(2 * time.Second)
	if isStandardIO == "true" && testingMode == "individualFiles" {
		// loop through the test cases
		for i := 0; i < numberOfTestCases; i++ {
			test_id := getEnv(fmt.Sprintf("TEST_%d_ID", i))
			// Get the input for the test case and write it to a file
			inputBase64 := getEnv(fmt.Sprintf("INPUT_%d_BASE64", i))
			input := decodeBase64(inputBase64)
			writeFile("console.in", input)

			// Run the code
			time.Sleep(1 * time.Second)
			err, memory, executionTime := runCPP("/executable", "console.in", "console.out", maxTime, maxMemory)
			if err != nil {
				fmt.Println("Error running executable")
				fmt.Println(err)
			}
			fmt.Println("Executable ran successfully")
			fmt.Println("Memory used:", memory)
			fmt.Println("Execution time:", executionTime)

			// Get the output for the test case and write it to a file
			outputBase64 := getEnv(fmt.Sprintf("OUTPUT_%d_BASE64", i))
			output := decodeBase64(outputBase64)
			writeFile("consoleRight.out", output)

			// Compare the output
			correctFile := compareFiles("console.out", "consoleRight.out")
			correctMemory := memory <= maxMemory
			correctTime := executionTime <= maxTime
			if correctFile && correctMemory && correctTime {
				fmt.Println("Test case", i, "passed")
				saveTestResult(test_id, true, memory, executionTime, "finished")
			} else {
				fmt.Println("Test case", i, "failed")
				saveTestResult(test_id, false, memory, executionTime, "finished")
			}

			fmt.Println("--------------------------------------------------")

			// Save the test result

			calculateScores(test_group_id, "running")

		}
	}
	calculateScores(test_group_id, "finished")

	time.Sleep(10 * time.Second)
}
