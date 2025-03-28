package terminal

import (
	"os"
	"strings"
	"testing"
)

// TestGenerateTestScript tests that the test script generation works properly
func TestGenerateTestScript(t *testing.T) {
	// Get the script content
	scriptContent, err := GenerateTestScript()
	if err != nil {
		t.Fatalf("Error generating test script: %v", err)
	}
	
	// Verify script content has expected sections
	expectedPhrases := []string{
		"TERMINAL_DEV_API_TOKEN",
		"TF_VAR_use_dev_environment=true",
		"provider \"terminal-coffee\"",
		"resource \"terminal_address\"",
		"resource \"terminal_payment_card\"",
		"resource \"terminal_coffee_order\"",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(scriptContent, phrase) {
			t.Errorf("Expected script to contain phrase: %s", phrase)
		}
	}
	
	// Write the script to a temporary file for testing
	tempFile := "./test_script_temp.sh"
	err = os.WriteFile(tempFile, []byte(scriptContent), 0755)
	if err != nil {
		t.Logf("Warning: Could not create test script file: %v", err)
		return
	}
	defer os.Remove(tempFile) // Clean up
	
	// Verify the file exists and is executable
	fileInfo, err := os.Stat(tempFile)
	if err != nil {
		t.Errorf("Error checking test script file: %v", err)
	} else if fileInfo.Mode().Perm()&0111 == 0 {
		t.Errorf("Test script is not executable")
	}
}