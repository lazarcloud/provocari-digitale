package main

import (
	"fmt"
	"os"
)

func finishExecution(error error) {
	// Send the execution time to the parent process
	// The parent process will read the time and print it
	if error != nil {
		fmt.Println("Error:", error)
		os.Exit(1)
	}
}
