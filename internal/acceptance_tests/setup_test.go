// Copyright (c) 2025 - Stacklet, Inc.

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
	if os.Getenv("TF_ACC") == "" || os.Getenv("TF_ACC_RECORD") != "" {
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
