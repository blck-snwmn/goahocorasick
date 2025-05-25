package goahocorasick

import (
	"strings"
	"testing"
)

func TestHighlyRepetitivePatterns(t *testing.T) {
	matcher := New()
	patterns := []string{
		"aaaaaaa",
		"aaaaaa",
		"aaaaa",
		"aaaa",
		"aaa",
		"aa",
		"a",
	}
	matcher.Build(patterns)
	
	text := strings.Repeat("a", 100)
	matches := matcher.FindAll(text)
	
	if len(matches) == 0 {
		t.Errorf("Expected many matches for repetitive patterns, got none")
	}
}

func TestPatternsSharingPrefixes(t *testing.T) {
	matcher := New()
	patterns := []string{
		"prefix",
		"prefixA",
		"prefixB",
		"prefixAB",
		"prefixABC",
		"prefixBCD",
	}
	matcher.Build(patterns)
	
	text := "prefixABCDEF"
	matches := matcher.FindAll(text)
	
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
	
	matcher := New()
	matcher.Build(patterns)
	
	text := "This is a test with pattern5 and pattern7"
	matches := matcher.FindAll(text)
	
	if len(matches) == 0 {
		t.Errorf("Expected matches for large pattern set, got none")
	}
}

func TestWorstCaseBacktracking(t *testing.T) {
	matcher := New()
	patterns := []string{"aaab", "aab", "ab", "b"}
	matcher.Build(patterns)
	
	text := strings.Repeat("a", 100) + "b"
	matches := matcher.FindAll(text)
	
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
	
	matcher := New()
	matcher.Build(patterns)
	
	text := strings.Repeat("a", 50) + "b" + strings.Repeat("a", 50)
	matches := matcher.FindAll(text)
	
	if len(matches) == 0 {
		t.Errorf("Expected matches for character bias patterns, got none")
	}
}