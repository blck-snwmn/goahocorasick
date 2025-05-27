package goahocorasick

import (
	"reflect"
	"testing"
)

func TestPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		expected []Match
	}{
		{
			name:     "basic matching",
			patterns: []string{"he", "she", "his", "hers"},
			text:     "ushers",
			expected: []Match{
				{Pattern: "she", Index: 1, Start: 1, End: 4},
				{Pattern: "he", Index: 0, Start: 2, End: 4},
				{Pattern: "hers", Index: 3, Start: 2, End: 6},
			},
		},
		{
			name:     "overlapping patterns",
			patterns: []string{"a", "ab", "abc", "bc", "c"},
			text:     "abc",
			expected: []Match{
				{Pattern: "a", Index: 0, Start: 0, End: 1},
				{Pattern: "ab", Index: 1, Start: 0, End: 2},
				{Pattern: "abc", Index: 2, Start: 0, End: 3},
				{Pattern: "bc", Index: 3, Start: 1, End: 3},
				{Pattern: "c", Index: 4, Start: 2, End: 3},
			},
		},
		{
			name:     "no matches",
			patterns: []string{"foo", "bar", "baz"},
			text:     "hello world",
			expected: []Match{},
		},
		{
			name:     "empty text",
			patterns: []string{"test"},
			text:     "",
			expected: []Match{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matcher := New()
			if err := matcher.Build(tc.patterns); err != nil {
				t.Fatalf("Build failed: %v", err)
			}

			matches, err := matcher.FindAll(tc.text)
			if err != nil {
				t.Fatalf("FindAll failed: %v", err)
			}

			if !reflect.DeepEqual(matches, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, matches)
			}
		})
	}
}

func TestUnicodeAndSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		text     string
		expected []Match
	}{
		{
			name:     "unicode patterns",
			patterns: []string{"こんにちは", "世界", "日本"},
			text:     "こんにちは世界、日本へようこそ",
			expected: []Match{
				{Pattern: "こんにちは", Index: 0, Start: 0, End: 5},
				{Pattern: "世界", Index: 1, Start: 5, End: 7},
				{Pattern: "日本", Index: 2, Start: 8, End: 10},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matcher := New()
			if err := matcher.Build(tc.patterns); err != nil {
				t.Fatalf("Build failed: %v", err)
			}

			matches, err := matcher.FindAll(tc.text)
			if err != nil {
				t.Fatalf("FindAll failed: %v", err)
			}

			if !reflect.DeepEqual(matches, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, matches)
			}
		})
	}
}


func TestEmptyPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{}
	err := matcher.Build(patterns)
	if err == nil {
		t.Fatal("Expected error for empty patterns, got nil")
	}
	
	text := "test"
	matches, err := matcher.FindAll(text)
	if err == nil {
		t.Fatal("Expected error when FindAll called without successful Build")
	}
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches for empty patterns, got %v", matches)
	}
}

func TestEmptyStringPattern(t *testing.T) {
	matcher := New()
	patterns := []string{"", "test"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "test"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "test", Index: 0, Start: 0, End: 4},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestDuplicatePatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"abc", "abc", "abc"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "abc"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "abc", Index: 0, Start: 0, End: 3},
		{Pattern: "abc", Index: 1, Start: 0, End: 3},
		{Pattern: "abc", Index: 2, Start: 0, End: 3},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestSingleCharacterPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"a", "b", "c", "d"}
	if err := matcher.Build(patterns); err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	text := "abcdcba"
	matches, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	
	expected := []Match{
		{Pattern: "a", Index: 0, Start: 0, End: 1},
		{Pattern: "b", Index: 1, Start: 1, End: 2},
		{Pattern: "c", Index: 2, Start: 2, End: 3},
		{Pattern: "d", Index: 3, Start: 3, End: 4},
		{Pattern: "c", Index: 2, Start: 4, End: 5},
		{Pattern: "b", Index: 1, Start: 5, End: 6},
		{Pattern: "a", Index: 0, Start: 6, End: 7},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestFindAllBeforeBuild(t *testing.T) {
	matcher := New()
	text := "test"
	matches, err := matcher.FindAll(text)
	
	if err == nil {
		t.Error("Expected error when FindAll called before Build, got nil")
	}
	if len(matches) != 0 {
		t.Errorf("Expected no matches before Build, got %v", matches)
	}
}

func TestRebuild(t *testing.T) {
	matcher := New()
	
	patterns1 := []string{"abc", "def"}
	if err := matcher.Build(patterns1); err != nil {
		t.Fatalf("First Build failed: %v", err)
	}
	
	text := "abcdef"
	matches1, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("First FindAll failed: %v", err)
	}
	
	expected1 := []Match{
		{Pattern: "abc", Index: 0, Start: 0, End: 3},
		{Pattern: "def", Index: 1, Start: 3, End: 6},
	}
	
	if !reflect.DeepEqual(matches1, expected1) {
		t.Errorf("First build: Expected %v, got %v", expected1, matches1)
	}
	
	patterns2 := []string{"xyz"}
	if err := matcher.Build(patterns2); err != nil {
		t.Fatalf("Second Build failed: %v", err)
	}
	
	matches2, err := matcher.FindAll(text)
	if err != nil {
		t.Fatalf("Second FindAll failed: %v", err)
	}
	
	if len(matches2) != 0 {
		t.Errorf("After rebuild: Expected no matches, got %v", matches2)
	}
}