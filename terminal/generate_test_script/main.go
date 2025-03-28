package main

import (
	"fmt"
	"os"

	"github.com/OZCAP/terraform-provider-terminal-coffee/terminal"
)

func main() {
	// Generate the test script
	scriptContent, err := terminal.GenerateTestScript()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating test script: %v\n", err)
		os.Exit(1)
	}

	// Write the script to a file
	err = os.WriteFile("test_dev_env.sh", []byte(scriptContent), 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing test script: %v\n", err)
		os.Exit(1)
	}

	// Print instructions
	fmt.Println("\nGenerated test script: test_dev_env.sh")
	fmt.Println("Edit this script with your development credentials and run to test the provider")
	fmt.Println(terminal.GetTestEnvironmentInstructions())
}