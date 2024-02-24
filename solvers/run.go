package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func runCPP(executableFile string, inputFile string) (err error, memory int64, executionTime time.Duration) {
	cmd := exec.Command(executableFile)

	// Open the input file
	infile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening input file")
		return err, 0, 0
	}
	defer infile.Close()

	// Redirect standard input from the file
	cmd.Stdin = infile

	// Create the output file
	outfile, err := os.Create("console.out")
	if err != nil {
		fmt.Println("Error creating output file")
		return err, 0, 0
	}
	defer outfile.Close()

	// Redirect standard output and standard error to the file
	cmd.Stdout = outfile
	cmd.Stderr = outfile

	startTime := time.Now() // Record the start time
	err = cmd.Start()
	if err != nil {
		return err, 0, 0
	}
	pid := cmd.Process.Pid
	fmt.Println("Executable running with PID:", pid)

	// Wait for the executable to finish
	err = cmd.Wait()
	executionTime = time.Since(startTime) // Calculate the execution time
	if err != nil {
		return err, 0, 0
	}
	fmt.Println("Executable finished successfully in", executionTime)

	// This is platform specific and may vary depending on the OS
	// Here we use the syscall package to get the resource usage of the process
	// See https://pkg.go.dev/syscall#Rusage for details
	var rusage syscall.Rusage
	err = syscall.Getrusage(syscall.RUSAGE_CHILDREN, &rusage)
	if err != nil {
		fmt.Println("Error getting resource usage")
		fmt.Println(err)
		return
	}
	// The Maxrss field is the maximum resident set size in kilobytes
	// This is the amount of memory occupied by the process in RAM
	// See https://man7.org/linux/man-pages/man2/getrusage.2.html for details
	memUsage := rusage.Maxrss
	return nil, memUsage, executionTime
}
