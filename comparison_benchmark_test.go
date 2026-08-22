package goahocorasick

import (
	"regexp"
	"strings"
	"testing"
)

func naiveSearch(text string, patterns []string) []Match {
	matches := make([]Match, 0)
	for idx, pattern := range patterns {
		start := 0
		for {
			pos := strings.Index(text[start:], pattern)
			if pos == -1 {
				break
			}
			actualPos := start + pos
			matches = append(matches, Match{
				Pattern: pattern,
				Index:   idx,
				Start:   actualPos,
				End:     actualPos + len(pattern),
			})
			start = actualPos + 1
		}
	}
	return matches
}

func generatePatterns(count int, length int) []string {
	patterns := make([]string, count)
	for i := 0; i < count; i++ {
		pattern := ""
		for j := 0; j < length; j++ {
			pattern += string(rune('a' + (i+j)%26))
		}
		patterns[i] = pattern
	}
	return patterns
}

func generateText(size int) string {
	return strings.Repeat("The quick brown fox jumps over the lazy dog. ", size)
}

func BenchmarkComparison5Patterns(b *testing.B) {
	patterns := []string{"quick", "brown", "fox", "lazy", "dog"}
	text := generateText(100)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparison20Patterns(b *testing.B) {
	patterns := generatePatterns(20, 4)
	text := generateText(100)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparison100Patterns(b *testing.B) {
	patterns := generatePatterns(100, 4)
	text := generateText(100)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		// Limit patterns for regexp to avoid "expression too large" error
		limitedPatterns := patterns[:20]
		pattern := strings.Join(limitedPatterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparisonShortText(b *testing.B) {
	patterns := generatePatterns(10, 3)
	text := generateText(10)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparisonLongText(b *testing.B) {
	patterns := generatePatterns(10, 3)
	text := generateText(1000)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparisonVeryLongText(b *testing.B) {
	patterns := generatePatterns(10, 3)
	text := generateText(10000)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile(pattern)
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}

func BenchmarkComparisonRealWorldPatterns(b *testing.B) {
	patterns := []string{
		"error", "warning", "info", "debug", "fatal",
		"exception", "failed", "success", "complete", "timeout",
		"connection", "request", "response", "server", "client",
	}
	text := strings.Repeat("2023-01-01 10:00:00 [ERROR] Failed to connect to server: connection timeout. Request failed with exception. ", 100)

	b.Run("AhoCorasick", func(b *testing.B) {
		builder := NewBuilder()
		builder.AddPatterns(patterns)
		matcher, err := builder.Build()
		if err != nil {
			b.Fatalf("Build failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("NaiveSearch", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = naiveSearch(text, patterns)
		}
	})

	b.Run("StringsContains", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			matches := 0
			for _, pattern := range patterns {
				if strings.Contains(text, pattern) {
					matches++
				}
			}
		}
	})

	b.Run("Regexp", func(b *testing.B) {
		pattern := strings.Join(patterns, "|")
		re, err := regexp.Compile("(?i)" + pattern) // Case insensitive
		if err != nil {
			b.Fatalf("Regexp compile failed: %v", err)
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = re.FindAllStringIndex(text, -1)
		}
	})
}
