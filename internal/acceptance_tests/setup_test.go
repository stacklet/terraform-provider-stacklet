// Copyright Stacklet, Inc. 2025, 2026

package acceptance_tests

import (
	"os"
	"testing"
)

// TestMain sets up/teardown the test environment around the test session.
func TestMain(m *testing.M) {
	ensureVars()
	os.Exit(m.Run())
}

func ensureVars() {
	// make it easier to get tests with paginated responses
	mustSetenv("STACKLET_PAGE_SIZE", "1")

	// Only set authentication-related vars when running tests in replay mode.
	if os.Getenv("TF_ACC") == "" || testMode() != TestModeReplay {
		return
	}

	mustSetenv("STACKLET_ENDPOINT", "https://fake")
	mustSetenv("STACKLET_API_KEY", "fake")
}

func mustSetenv(name string, value string) {
	if err := os.Setenv(name, value); err != nil {
		panic(err)
	}
}
