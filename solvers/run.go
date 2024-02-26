package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func runCPP(executableFile string, inputFile string, outputFile string, maxTime time.Duration, memoryLimit int64) (err error, memory int64, executionTime time.Duration) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), maxTime)
	defer cancel()

	cmd := exec.CommandContext(ctx, executableFile)

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
	outfile, err := os.Create(outputFile)
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

	// Channel to communicate if execution is completed
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Timeout reached
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("Error killing process:", err)
		}
		fmt.Println("Execution timed out")
		return ctx.Err(), 0, time.Since(startTime)
	case err := <-done:
		// Execution completed
		if err != nil {
			return err, 0, 0
		}
		fmt.Println("Executable finished successfully")
		executionTime = time.Since(startTime)
	}

	// Get memory usage
	var rusage syscall.Rusage
	err = syscall.Getrusage(syscall.RUSAGE_CHILDREN, &rusage)
	if err != nil {
		fmt.Println("Error getting resource usage")
		return err, 0, executionTime
	}
	// The Maxrss field is the maximum resident set size in kilobytes
	memUsage := rusage.Maxrss
	return nil, memUsage, executionTime
}
