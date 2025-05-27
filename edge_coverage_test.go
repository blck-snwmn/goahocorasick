package goahocorasick

import (
	"testing"
)

func TestSelfReferencePreventionEdgeCase(t *testing.T) {
	matcher := New()
	patterns := []string{"a"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "a"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) != 1 || matches[0].Pattern != "a" {
		t.Errorf("Expected one match for 'a', got %v", matches)
	}
}

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

func TestSinglePatternNoFailureLinks(t *testing.T) {
	matcher := New()
	patterns := []string{"single"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "single"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) != 1 || matches[0].Pattern != "single" {
		t.Errorf("Expected one match for 'single', got %v", matches)
	}
}