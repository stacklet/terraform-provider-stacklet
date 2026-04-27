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
	// Only set vars when running tests in replay mode.
	if os.Getenv("TF_ACC") == "" || testMode() != TestModeReplay {
		return
	}

	vars := map[string]string{
		"STACKLET_ENDPOINT": "https://fake",
		"STACKLET_API_KEY":  "fake",
	}
	for name, value := range vars {
		if err := os.Setenv(name, value); err != nil {
			panic(err)
		}
	}
}
