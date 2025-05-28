# goahocorasick

A Go implementation of the Aho-Corasick algorithm for efficient multiple pattern matching.

## Description

goahocorasick provides a fast string matching algorithm that can find all occurrences of multiple patterns in a text with a single pass. The algorithm constructs a finite state machine from a set of patterns and matches them in O(n) time complexity where n is the length of the text.

## Installation

```bash
go get github.com/blck-snwmn/goahocorasick
```

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/blck-snwmn/goahocorasick"
)

func main() {
    // Create a new builder
    builder := goahocorasick.NewBuilder()
    
    // Add patterns
    builder.AddPatterns([]string{"he", "she", "his", "hers"})
    
    // Build the matcher
    matcher, err := builder.Build()
    if err != nil {
        log.Fatal(err)
    }
    
    // Find all matches
    matches, err := matcher.FindAll("ushers")
    if err != nil {
        log.Fatal(err)
    }
    
    // Print results
    for _, match := range matches {
        fmt.Printf("Found '%s' at position %d-%d\n", match.Pattern, match.Start, match.End)
    }
}
```

## Features

- Thread-safe matcher (immutable after building)
- UTF-8 support
- Optimized for ASCII characters
- Builder pattern for easy construction