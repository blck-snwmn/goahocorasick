package goahocorasick

import (
	"testing"
)

func TestSelfReferencePreventionEdgeCase(t *testing.T) {
	matcher := New()
	patterns := []string{"a"}
	matcher.Build(patterns)
	
	text := "a"
	matches := matcher.FindAll(text)
	
	if len(matches) != 1 || matches[0].Pattern != "a" {
		t.Errorf("Expected one match for 'a', got %v", matches)
	}
}

func TestEmptyOutputSuffixLink(t *testing.T) {
	matcher := New()
	patterns := []string{"xyz", "yz"}
	matcher.Build(patterns)
	
	text := "xyz"
	matches := matcher.FindAll(text)
	
	expectedCount := 2
	if len(matches) != expectedCount {
		t.Errorf("Expected %d matches, got %d", expectedCount, len(matches))
	}
}

func TestSinglePatternNoFailureLinks(t *testing.T) {
	matcher := New()
	patterns := []string{"single"}
	matcher.Build(patterns)
	
	text := "single"
	matches := matcher.FindAll(text)
	
	if len(matches) != 1 || matches[0].Pattern != "single" {
		t.Errorf("Expected one match for 'single', got %v", matches)
	}
}