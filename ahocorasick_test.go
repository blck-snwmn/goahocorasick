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