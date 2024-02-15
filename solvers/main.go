package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cppSourceBase64 := os.Getenv("CPP_SOURCE_BASE64")
	if cppSourceBase64 == "" {
		fmt.Println("No source code provided")
		return
	}
	cppSource, err := base64.StdEncoding.DecodeString(cppSourceBase64)
	if err != nil {
		fmt.Println("Error decoding base64")
		return
	}
	fmt.Println("Source code decoded successfully")
	err = os.WriteFile("main.cpp", cppSource, 0755)
	if err != nil {
		fmt.Println("Error writing source code")
		return
	}
	defer os.Remove("main.cpp")

	fmt.Println("Source code written successfully")

	// Compile the source code
	cmd := exec.Command("g++", "main.cpp", "-o", "executable")
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error compiling source code")
		fmt.Println(err)
		return
	}
	fmt.Println("Source code compiled successfully")

	// Run the executable and get its process ID
	cmd = exec.Command("./executable")
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error running executable")
		fmt.Println(err)
		return
	}
	pid := cmd.Process.Pid
	fmt.Println("Executable running with PID:", pid)

	// Wait for the executable to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for executable")
		fmt.Println(err)
		return
	}
	fmt.Println("Executable finished successfully")

	// Get the memory usage of the executable from the OS
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
	fmt.Println("Memory usage:", memUsage, "KB")

	// Clean up temporary file
	os.Remove("executable")
}
