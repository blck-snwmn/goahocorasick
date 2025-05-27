package goahocorasick

import (
	"testing"
)

// Test various pattern counts to see performance impact
func BenchmarkSinglePatternPerformance(b *testing.B) {
	pattern := "error"
	text := "This is an error message with error codes and error handling"
	
	builder := NewBuilder()
	builder.AddPattern(pattern)
	matcher, err := builder.Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matcher.FindAll(text)
	}
}

func BenchmarkThreePatternsPerformance(b *testing.B) {
	patterns := []string{"error", "warning", "info"}
	text := "2023-01-01 [ERROR] System error. [WARNING] Low memory. [INFO] Process started."
	
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
}

func BenchmarkFindAllOverhead(b *testing.B) {
	// Minimal case to measure FindAll overhead
	pattern := "x"
	text := "abcdefghijklmnopqrstuvwxyz"
	
	builder := NewBuilder()
	builder.AddPattern(pattern)
	matcher, err := builder.Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matcher.FindAll(text)
	}
}

func BenchmarkNoMatchesOverhead(b *testing.B) {
	// Case where no patterns match
	pattern := "xyz"
	text := "abcdefghijklmnopqrstuvw"
	
	builder := NewBuilder()
	builder.AddPattern(pattern)
	matcher, err := builder.Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}
	
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = matcher.FindAll(text)
	}
}