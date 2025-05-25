package goahocorasick

import (
	"reflect"
	"testing"
)

func TestComplexFailureLinks(t *testing.T) {
	matcher := New()
	patterns := []string{"abcde", "cde", "bcde", "de", "e"}
	matcher.Build(patterns)
	
	text := "abcde"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "abcde", Index: 0, Start: 0, End: 5},
		{Pattern: "bcde", Index: 2, Start: 1, End: 5},
		{Pattern: "cde", Index: 1, Start: 2, End: 5},
		{Pattern: "de", Index: 3, Start: 3, End: 5},
		{Pattern: "e", Index: 4, Start: 4, End: 5},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestDeepNestingPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"aaaa", "aaa", "aa", "a"}
	matcher.Build(patterns)
	
	text := "aaaa"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "a", Index: 3, Start: 0, End: 1},
		{Pattern: "aa", Index: 2, Start: 0, End: 2},
		{Pattern: "a", Index: 3, Start: 1, End: 2},
		{Pattern: "aaa", Index: 1, Start: 0, End: 3},
		{Pattern: "aa", Index: 2, Start: 1, End: 3},
		{Pattern: "a", Index: 3, Start: 2, End: 3},
		{Pattern: "aaaa", Index: 0, Start: 0, End: 4},
		{Pattern: "aaa", Index: 1, Start: 1, End: 4},
		{Pattern: "aa", Index: 2, Start: 2, End: 4},
		{Pattern: "a", Index: 3, Start: 3, End: 4},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestSuffixLinkOutput(t *testing.T) {
	matcher := New()
	patterns := []string{"she", "he", "hers"}
	matcher.Build(patterns)
	
	text := "shers"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "she", Index: 0, Start: 0, End: 3},
		{Pattern: "he", Index: 1, Start: 1, End: 3},
		{Pattern: "hers", Index: 2, Start: 1, End: 5},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestFailureFunctionWithOverlap(t *testing.T) {
	matcher := New()
	patterns := []string{"abab", "bab", "ab"}
	matcher.Build(patterns)
	
	text := "abababab"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "ab", Index: 2, Start: 0, End: 2},
		{Pattern: "abab", Index: 0, Start: 0, End: 4},
		{Pattern: "bab", Index: 1, Start: 1, End: 4},
		{Pattern: "ab", Index: 2, Start: 2, End: 4},
		{Pattern: "abab", Index: 0, Start: 2, End: 6},
		{Pattern: "bab", Index: 1, Start: 3, End: 6},
		{Pattern: "ab", Index: 2, Start: 4, End: 6},
		{Pattern: "abab", Index: 0, Start: 4, End: 8},
		{Pattern: "bab", Index: 1, Start: 5, End: 8},
		{Pattern: "ab", Index: 2, Start: 6, End: 8},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}

func TestPrefixSuffixPatterns(t *testing.T) {
	matcher := New()
	patterns := []string{"abcab", "cab", "ab", "bcab"}
	matcher.Build(patterns)
	
	text := "abcabcab"
	matches := matcher.FindAll(text)
	
	expected := []Match{
		{Pattern: "ab", Index: 2, Start: 0, End: 2},
		{Pattern: "abcab", Index: 0, Start: 0, End: 5},
		{Pattern: "bcab", Index: 3, Start: 1, End: 5},
		{Pattern: "cab", Index: 1, Start: 2, End: 5},
		{Pattern: "ab", Index: 2, Start: 3, End: 5},
		{Pattern: "abcab", Index: 0, Start: 3, End: 8},
		{Pattern: "bcab", Index: 3, Start: 4, End: 8},
		{Pattern: "cab", Index: 1, Start: 5, End: 8},
		{Pattern: "ab", Index: 2, Start: 6, End: 8},
	}
	
	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("Expected %v, got %v", expected, matches)
	}
}