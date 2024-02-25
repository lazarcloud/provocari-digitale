package main

import (
	"fmt"
	"time"
)

var validTestingModes = []string{"individualFiles"}

func main() {
	// Get problem details

	isStandardIO := getEnv("USES_STANDARD_IO")
	testingMode := getEnv("TEST_MODE")
	problemInputFile := getEnv("PROBLEM_INPUT_FILE")
	problemOutputFile := getEnv("PROBLEM_OUTPUT_FILE")
	numberOfTestCases := parseInt(getEnv("NUMBER_OF_TEST_CASES"))
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

	// Compile the source code
	compileCPP("./main.cpp", "executable")
	fmt.Println("Source code compiled successfully")

	if isStandardIO == "true" && testingMode == "individualFiles" {
		// loop through the test cases
		for i := 0; i < numberOfTestCases; i++ {
			test_id := getEnv(fmt.Sprintf("TEST_%d_ID", i))
			// Get the input for the test case and write it to a file
			inputBase64 := getEnv(fmt.Sprintf("INPUT_%d_BASE64", i))
			input := decodeBase64(inputBase64)
			writeFile("console.in", input)

			// Run the code
			err, memory, executionTime := runCPP("/executable", "console.in", "console.out")
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
			correct := compareFiles("console.out", "consoleRight.out")
			if correct {
				fmt.Println("Test case", i, "passed")
			} else {
				fmt.Println("Test case", i, "failed")
			}

			fmt.Println("--------------------------------------------------")

			// Save the test result
			saveTestResult(test_id, correct, memory, executionTime)
		}
	}

	time.Sleep(10 * time.Second)
}
