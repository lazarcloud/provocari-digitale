package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
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
	fmt.Println("Source code written successfully")
	//compile the source code
	cmd := exec.Command("g++", "main.cpp", "-o", "executable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error compiling source code")
		return
	}
	fmt.Println("Source code compiled successfully")

	//measure the time and memory usage
	var m1, m2 runtime.MemStats
	start := time.Now()
	runtime.ReadMemStats(&m1) //read memory stats before execution
	cmd = exec.Command("./executable")
	output, err = cmd.CombinedOutput() //run the executable
	runtime.ReadMemStats(&m2)          //read memory stats after execution
	elapsed := time.Since(start)       //calculate elapsed time

	if err != nil {
		fmt.Println("Error running program")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Max Memory:", m2.TotalAlloc-m1.TotalAlloc) //print the difference in total allocated memory
	fmt.Println("Elapsed Time:", elapsed)                   //print the elapsed time
	fmt.Println("Output:", string(output))                  //print the output of the executable
	return

	// Clean up temporary file
	os.Remove("executable")
}
