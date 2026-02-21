package main

import (
	"os"
	"testing"
)

// TestMain wraps the entire test suite to protect sensitive files from being
// overwritten by tests that trigger UI actions (e.g. environment switching,
// which calls SaveConfig → writes o8n-env.yaml).
func TestMain(m *testing.M) {
	const envFile = "o8n-env.yaml"

	// Back up the current env file before any test runs.
	original, readErr := os.ReadFile(envFile)

	code := m.Run()

	// Restore original content regardless of test outcome.
	if readErr == nil {
		_ = os.WriteFile(envFile, original, 0600)
	}

	os.Exit(code)
}
