package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// Prepare the go test command
	// You can specify a particular package or file by adding arguments to the command
	cmd := exec.Command("go", "test", "github.com/askiada/go-sql-test/internal/parser", "-run", "TestRunSQL")

	// Run the command and capture the output and error
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output)) //nolint:forbidigo // Print the output before exiting

	if err != nil {
		// Check if the error is an ExitError
		if exitErr, ok := err.(*exec.ExitError); ok { //nolint:errorlint // error is never wrapped
			// Check if the underlying system call error is available and extract the exit code
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}

		os.Exit(1) // Default exit code for general errors
	}
}
