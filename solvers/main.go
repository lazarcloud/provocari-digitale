package main

import (
	"fmt"
	"time"
)

func main() {
	// Get problem details
	problemInputFile := getEnv("PROBLEM_INPUT_FILE")
	problemOutputFile := getEnv("PROBLEM_OUTPUT_FILE")
	numberOfTestCases := parseInt(getEnv("NUMBER_OF_TEST_CASES"))
	// print the problem details
	fmt.Println("Problem input file:", problemInputFile)
	fmt.Println("Problem output file:", problemOutputFile)
	fmt.Println("Number of test cases:", numberOfTestCases)

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

	// loop through the test cases
	for i := 0; i < numberOfTestCases; i++ {
		// Get the input for the test case and write it to a file
		inputBase64 := getEnv(fmt.Sprintf("INPUT_%d_BASE64", i))
		input := decodeBase64(inputBase64)
		writeFile(problemInputFile, input)

		// Run the code
		err, memory, executionTime := runCPP("/executable", problemInputFile)
		if err != nil {
			fmt.Println("Error running executable")
			fmt.Println(err)
		}
		fmt.Println("Executable ran successfully")
		fmt.Println("Memory used:", memory)
		fmt.Println("Execution time:", executionTime)
	}

	time.Sleep(10 * time.Second)
}
