package goahocorasick

import (
	"strings"
	"testing"
)


func TestPatternsSharingPrefixes(t *testing.T) {
	builder := NewBuilder()
	patterns := []string{
		"prefix",
		"prefixA",
		"prefixB",
		"prefixAB",
		"prefixABC",
		"prefixBCD",
	}
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "prefixABCDEF"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expectedPatterns := map[string]bool{
		"prefix":    true,
		"prefixA":   true,
		"prefixAB":  true,
		"prefixABC": true,
	}
	
	for _, match := range matches {
		if !expectedPatterns[match.Pattern] {
			t.Errorf("Unexpected match: %v", match.Pattern)
		}
		delete(expectedPatterns, match.Pattern)
	}
	
	if len(expectedPatterns) > 0 {
		t.Errorf("Missing expected patterns: %v", expectedPatterns)
	}
}

func TestManyPatternsMemoryUsage(t *testing.T) {
	patterns := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		patterns[i] = "pattern" + string(rune('0'+i%10))
	}
	
	builder := NewBuilder()
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "This is a test with pattern5 and pattern7"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) == 0 {
		t.Errorf("Expected matches for large pattern set, got none")
	}
}

func TestWorstCaseBacktracking(t *testing.T) {
	builder := NewBuilder()
	patterns := []string{"aaab", "aab", "ab", "b"}
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := strings.Repeat("a", 100) + "b"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	foundB := false
	for _, match := range matches {
		if match.Pattern == "b" && match.Start == 100 {
			foundB = true
			break
		}
	}
	
	if !foundB {
		t.Errorf("Expected to find 'b' at position 100")
	}
}

func TestCharacterBiasPatterns(t *testing.T) {
	patterns := make([]string, 100)
	for i := 0; i < 100; i++ {
		length := i%10 + 1
		patterns[i] = strings.Repeat("a", length) + "b"
	}
	
	builder := NewBuilder()
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := strings.Repeat("a", 50) + "b" + strings.Repeat("a", 50)
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) == 0 {
		t.Errorf("Expected matches for character bias patterns, got none")
	}
}