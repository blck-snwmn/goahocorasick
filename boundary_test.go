package goahocorasick

import (
	"reflect"
	"strings"
	"testing"
)

func TestPatternLongerThanText(t *testing.T) {
	matcher := New()
	patterns := []string{"abcdef"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "abc"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches when pattern is longer than text, got %v", matches)
	}
}

func TestSingleCharacterText(t *testing.T) {
	matcher := New()
	patterns := []string{"a", "ab", "b"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "a"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "a", Index: 0, Start: 0, End: 1},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestPatternEqualsText(t *testing.T) {
	matcher := New()
	patterns := []string{"exact"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "exact"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "exact", Index: 0, Start: 0, End: 5},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestMatchAtTextBoundaries(t *testing.T) {
	matcher := New()
	patterns := []string{"start", "end", "middle"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "start in the middle at the end"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "start", Index: 0, Start: 0, End: 5},
		{Pattern: "middle", Index: 2, Start: 13, End: 19},
		{Pattern: "end", Index: 1, Start: 27, End: 30},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestVeryLongPattern(t *testing.T) {
	longPattern := strings.Repeat("a", 1000)
	matcher := New()
	patterns := []string{longPattern, "b"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := longPattern + "b"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: longPattern, Index: 0, Start: 0, End: 1000},
		{Pattern: "b", Index: 1, Start: 1000, End: 1001},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v matches, got %v", len(expected), len(matches))
	}
}

func TestManyShortPatterns(t *testing.T) {
	patterns := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		patterns[i] = string(rune('a' + i%26))
	}
	
	matcher := New()
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "abcdefghijklmnopqrstuvwxyz"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	if len(matches) < 26 {
		t.Errorf("Expected at least 26 matches, got %v", len(matches))
	}
}

func TestControlCharacters(t *testing.T) {
	matcher := New()
	patterns := []string{"line\n", "\ttab", "return\r"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "line\ntext\ttabtext return\rtext"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "line\n", Index: 0, Start: 0, End: 5},
		{Pattern: "\ttab", Index: 1, Start: 9, End: 13},
		{Pattern: "return\r", Index: 2, Start: 18, End: 25},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestEmojisAndComplexUnicode(t *testing.T) {
	matcher := New()
	patterns := []string{"😀", "👨‍👩‍👧‍👦", "🇯🇵", "𠮷"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "Hello 😀 family 👨‍👩‍👧‍👦 from 🇯🇵 and 𠮷"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "😀", Index: 0, Start: 6, End: 7},
		{Pattern: "👨‍👩‍👧‍👦", Index: 1, Start: 15, End: 22},
		{Pattern: "🇯🇵", Index: 2, Start: 28, End: 30},
		{Pattern: "𠮷", Index: 3, Start: 35, End: 36},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}