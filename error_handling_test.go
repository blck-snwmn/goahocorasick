package goahocorasick

import (
	"testing"
)

func TestBuildErrors(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		wantErr  bool
	}{
		{
			name:     "nil patterns",
			patterns: nil,
			wantErr:  true,
		},
		{
			name:     "empty patterns slice",
			patterns: []string{},
			wantErr:  true,
		},
		{
			name:     "all empty patterns",
			patterns: []string{"", "", ""},
			wantErr:  true,
		},
		{
			name:     "valid patterns",
			patterns: []string{"test", "pattern"},
			wantErr:  false,
		},
		{
			name:     "mixed empty and valid patterns",
			patterns: []string{"", "valid", "", "pattern", ""},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			builder.AddPatterns(tt.patterns)
			_, err := builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindAllErrors(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Matcher
		text    string
		wantErr bool
	}{
		{
			name: "matcher not built",
			setup: func() *Matcher {
				// Return nil to simulate matcher not built
				return nil
			},
			text:    "test text",
			wantErr: true,
		},
		{
			name: "valid matcher with valid text",
			setup: func() *Matcher {
				builder := NewBuilder()
				builder.AddPatterns([]string{"test"})
				m, _ := builder.Build()
				return m
			},
			text:    "test text",
			wantErr: false,
		},
		{
			name: "valid matcher with invalid UTF-8",
			setup: func() *Matcher {
				builder := NewBuilder()
				builder.AddPatterns([]string{"test"})
				m, _ := builder.Build()
				return m
			},
			text:    "test \xbd text",
			wantErr: true,
		},
		{
			name: "valid matcher with empty text",
			setup: func() *Matcher {
				builder := NewBuilder()
				builder.AddPatterns([]string{"test"})
				m, _ := builder.Build()
				return m
			},
			text:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			_, err := m.FindAll(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildAfterBuildWithError(t *testing.T) {
	// First build with error
	builder1 := NewBuilder()
	builder1.AddPatterns(nil)
	_, err := builder1.Build()
	if err == nil {
		t.Error("Expected error for nil patterns")
	}

	// Second build with valid patterns should succeed
	builder2 := NewBuilder()
	builder2.AddPatterns([]string{"valid", "patterns"})
	m, err := builder2.Build()
	if err != nil {
		t.Errorf("Expected no error for valid patterns, got: %v", err)
	}

	// FindAll should work after successful build
	matches, err := m.FindAll("valid text with patterns")
	if err != nil {
		t.Errorf("Expected no error for FindAll after successful build, got: %v", err)
	}
	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}
}
