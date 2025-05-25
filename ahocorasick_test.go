package goahocorasick

import (
	"reflect"
	"testing"
)

func TestBasicMatching(t *testing.T) {
	matcher := New()
	patterns := []string{"he", "she", "his", "hers"}
	matcher.Build(patterns)
	
	text := "ushers"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "she", Index: 1, Start: 1, End: 4},
		{Pattern: "he", Index: 0, Start: 2, End: 4},
		{Pattern: "hers", Index: 3, Start: 2, End: 6},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestOverlappingPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"a", "ab", "abc", "bc", "c"}
	matcher.Build(patterns)
	
	text := "abc"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "a", Index: 0, Start: 0, End: 1},
		{Pattern: "ab", Index: 1, Start: 0, End: 2},
		{Pattern: "abc", Index: 2, Start: 0, End: 3},
		{Pattern: "bc", Index: 3, Start: 1, End: 3},
		{Pattern: "c", Index: 4, Start: 2, End: 3},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestNoMatches(t *testing.T) {
	matcher := New()
	patterns := []string{"foo", "bar", "baz"}
	matcher.Build(patterns)
	
	text := "hello world"
	matches := matcher.FindAll(text)
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches, got %v", matches)
	}
}

func TestEmptyText(t *testing.T) {
	matcher := New()
	patterns := []string{"test"}
	matcher.Build(patterns)
	
	text := ""
	matches := matcher.FindAll(text)
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches for empty text, got %v", matches)
	}
}

func TestUnicodePatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"こんにちは", "世界", "日本"}
	matcher.Build(patterns)
	
	text := "こんにちは世界、日本へようこそ"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "こんにちは", Index: 0, Start: 0, End: 5},
		{Pattern: "世界", Index: 1, Start: 5, End: 7},
		{Pattern: "日本", Index: 2, Start: 8, End: 10},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestRepeatedPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"aa", "a"}
	matcher.Build(patterns)
	
	text := "aaa"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "a", Index: 1, Start: 0, End: 1},
		{Pattern: "aa", Index: 0, Start: 0, End: 2},
		{Pattern: "a", Index: 1, Start: 1, End: 2},
		{Pattern: "aa", Index: 0, Start: 1, End: 3},
		{Pattern: "a", Index: 1, Start: 2, End: 3},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestEmptyPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{}
	matcher.Build(patterns)
	
	text := "test"
	matches := matcher.FindAll(text)
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches for empty patterns, got %v", matches)
	}
}

func TestEmptyStringPattern(t *testing.T) {
	matcher := New()
	patterns := []string{"", "test"}
	matcher.Build(patterns)
	
	text := "test"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "abc"
	matches := matcher.FindAll(text)
	
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
	matcher.Build(patterns)
	
	text := "abcdcba"
	matches := matcher.FindAll(text)
	
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
	matches := matcher.FindAll(text)
	
	if len(matches) != 0 {
		t.Errorf("Expected no matches before Build, got %v", matches)
	}
}

func TestRebuild(t *testing.T) {
	matcher := New()
	
	patterns1 := []string{"abc", "def"}
	matcher.Build(patterns1)
	
	text := "abcdef"
	matches1 := matcher.FindAll(text)
	
	expected1 := []Match{
		{Pattern: "abc", Index: 0, Start: 0, End: 3},
		{Pattern: "def", Index: 1, Start: 3, End: 6},
	}
	
	if !reflect.DeepEqual(matches1, expected1) {
		t.Errorf("First build: Expected %v, got %v", expected1, matches1)
	}
	
	patterns2 := []string{"xyz"}
	matcher.Build(patterns2)
	
	matches2 := matcher.FindAll(text)
	
	if len(matches2) != 0 {
		t.Errorf("After rebuild: Expected no matches, got %v", matches2)
	}
}