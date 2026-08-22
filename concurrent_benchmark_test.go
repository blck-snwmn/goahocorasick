package goahocorasick

import (
	"sync"
	"testing"
)

func BenchmarkConcurrentFindAll(b *testing.B) {
	patterns := []string{"error", "warning", "info", "debug", "fatal"}
	text := "2023-01-01 10:00:00 [ERROR] Failed to connect to server. [WARNING] Low memory. [INFO] Started."

	builder := NewBuilder()
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}

	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("Concurrent-2", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = matcher.FindAll(text)
			}
		})
	})

	b.Run("Concurrent-4", func(b *testing.B) {
		b.SetParallelism(4)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = matcher.FindAll(text)
			}
		})
	})

	b.Run("Concurrent-8", func(b *testing.B) {
		b.SetParallelism(8)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = matcher.FindAll(text)
			}
		})
	})
}

func BenchmarkConcurrentManyPatterns(b *testing.B) {
	// Create 100 patterns
	patterns := make([]string, 100)
	for i := 0; i < 100; i++ {
		patterns[i] = string(rune('a'+i%26)) + string(rune('a'+(i+1)%26)) + string(rune('a'+(i+2)%26))
	}
	text := "The quick brown fox jumps over the lazy dog. " +
		"Pack my box with five dozen liquor jugs. " +
		"How vexingly quick daft zebras jump! " +
		"The five boxing wizards jump quickly."

	builder := NewBuilder()
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		b.Fatalf("Build failed: %v", err)
	}

	b.Run("Sequential", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = matcher.FindAll(text)
		}
	})

	b.Run("Concurrent-8", func(b *testing.B) {
		b.SetParallelism(8)
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = matcher.FindAll(text)
			}
		})
	})
}

func TestConcurrentAccess(t *testing.T) {
	patterns := []string{"test", "pattern", "match"}
	text := "This is a test pattern to match."

	builder := NewBuilder()
	builder.AddPatterns(patterns)
	matcher, err := builder.Build()
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// Run concurrent FindAll operations
	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	errors := make(chan error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			matches, err := matcher.FindAll(text)
			if err != nil {
				errors <- err
				return
			}
			if len(matches) != 3 {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}
