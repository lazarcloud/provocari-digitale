package main

import (
	"fmt"
	"os/exec"
)

func compileCPP(sourceFile string, executableFile string) {
	cmd := exec.Command("g++", sourceFile, "-o", executableFile)
	_, err := cmd.CombinedOutput()
	if err != nil {
		finishExecution(fmt.Errorf("Error compiling source code %s", err.Error()))
	}
}
