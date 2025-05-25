package goahocorasick

import (
	"strings"
	"testing"
)

func BenchmarkSmallPatterns(b *testing.B) {
	matcher := New()
	patterns := []string{"he", "she", "his", "hers"}
	matcher.Build(patterns)
	
	text := strings.Repeat("ushers ", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = matcher.FindAll(text)
	}
}

func BenchmarkManyPatterns(b *testing.B) {
	matcher := New()
	patterns := make([]string, 100)
	for i := 0; i < 100; i++ {
		patterns[i] = string(rune('a' + i%26)) + string(rune('a' + (i+1)%26))
	}
	matcher.Build(patterns)
	
	text := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = matcher.FindAll(text)
	}
}

func BenchmarkLongText(b *testing.B) {
	matcher := New()
	patterns := []string{"pattern", "test", "benchmark", "performance"}
	matcher.Build(patterns)
	
	text := strings.Repeat("This is a test of the pattern matching performance benchmark. ", 1000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = matcher.FindAll(text)
	}
}

func BenchmarkUnicodePatterns(b *testing.B) {
	matcher := New()
	patterns := []string{"こんにちは", "世界", "日本", "東京", "大阪"}
	matcher.Build(patterns)
	
	text := strings.Repeat("こんにちは世界、日本の東京と大阪へようこそ。", 100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = matcher.FindAll(text)
	}
}

func BenchmarkBuildTime(b *testing.B) {
	patterns := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		patterns[i] = "pattern" + string(rune('0' + i%10))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher := New()
		matcher.Build(patterns)
	}
}