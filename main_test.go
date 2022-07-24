package main

import (
	"testing"
)

func TestSuccessCount(t *testing.T) {
	ts := TestSuite{
		Tests:    "10",
		Failures: "0",
		Errors:   "0",
		Skipped:  "0",
	}

	if ts.GetSuccessCount() != 10 {
		t.Errorf("Expected 10, got %d", ts.GetSuccessCount())
	}
}
