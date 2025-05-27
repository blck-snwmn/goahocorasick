package goahocorasick

import (
	"testing"
)


func TestEmptyOutputSuffixLink(t *testing.T) {
	matcher := New()
	patterns := []string{"xyz", "yz"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "xyz"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expectedCount := 2
	if len(matches) != expectedCount {
		t.Errorf("Expected %d matches, got %d", expectedCount, len(matches))
	}
}

