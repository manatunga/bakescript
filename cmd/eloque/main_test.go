package main

import "testing"

func TestStartupMessage(t *testing.T) {
	expected := "Eloque - programming language under development"
	result := startupMessage()

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
