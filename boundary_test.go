package goahocorasick

import (
	"reflect"
	"strings"
	"testing"
)

func TestPatternLongerThanText(t *testing.T) {
	matcher := New()
	patterns := []string{"abcdef"}
	matcher.Build(patterns)
	
	text := "abc"
	matches := matcher.FindAll(text)
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches when pattern is longer than text, got %v", matches)
	}
}

func TestSingleCharacterText(t *testing.T) {
	matcher := New()
	patterns := []string{"a", "ab", "b"}
	matcher.Build(patterns)
	
	text := "a"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "exact"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "start in the middle at the end"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := longPattern + "b"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "abcdefghijklmnopqrstuvwxyz"
	matches := matcher.FindAll(text)
	
	if len(matches) < 26 {
		t.Errorf("Expected at least 26 matches, got %v", len(matches))
	}
}

func TestControlCharacters(t *testing.T) {
	matcher := New()
	patterns := []string{"line\n", "\ttab", "return\r"}
	matcher.Build(patterns)
	
	text := "line\ntext\ttabtext return\rtext"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "Hello 😀 family 👨‍👩‍👧‍👦 from 🇯🇵 and 𠮷"
	matches := matcher.FindAll(text)
	
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