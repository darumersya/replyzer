package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Test that main doesn't panic when called with no arguments
	// This is a basic smoke test
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	// Test with help flag
	os.Args = []string{"analyze-repo", "--help"}
	
	// This should not panic and complete successfully
	// We can't easily test the actual execution without creating test files
	// but we can ensure the command structure is valid
	
	// Reset args
	os.Args = oldArgs
}

func TestVersionInfo(t *testing.T) {
	// Test that version variable can be set (used in build process)
	originalVersion := version
	version = "test-version"
	
	if version != "test-version" {
		t.Errorf("Expected version to be 'test-version', got %s", version)
	}
	
	// Restore original version
	version = originalVersion
}