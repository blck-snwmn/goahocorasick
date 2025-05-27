// Package goahocorasick implements the Aho-Corasick algorithm for efficient
// multiple pattern matching. The algorithm constructs a finite state machine
// from a set of patterns and can find all occurrences of any pattern in a text
// in linear time relative to the length of the text.
//
// Example usage:
//
//	builder := goahocorasick.NewBuilder()
//	builder.AddPatterns([]string{"he", "she", "his", "hers"})
//	matcher, err := builder.Build()
//	if err != nil {
//		log.Fatal(err)
//	}
//	matches, err := matcher.FindAll("ushers")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, match := range matches {
//		fmt.Printf("Found '%s' at position %d-%d\n", match.Pattern, match.Start, match.End)
//	}
package goahocorasick

import (
	"errors"
	"unicode/utf8"
)

// Builder is used to construct an immutable Matcher.
// The Builder pattern allows for a clear separation between the construction
// and usage phases, enabling lock-free concurrent access to the resulting Matcher.
type Builder struct {
	// patterns stores patterns to be built into the matcher
	patterns []string
}

// NewBuilder creates a new Builder instance for constructing a Matcher.
func NewBuilder() *Builder {
	return &Builder{
		patterns: make([]string, 0),
	}
}

// AddPattern adds a pattern to the builder.
// Empty patterns are ignored.
func (b *Builder) AddPattern(pattern string) *Builder {
	if len(pattern) > 0 {
		b.patterns = append(b.patterns, pattern)
	}
	return b
}

// AddPatterns adds multiple patterns to the builder.
// Empty patterns are ignored.
func (b *Builder) AddPatterns(patterns []string) *Builder {
	for _, pattern := range patterns {
		b.AddPattern(pattern)
	}
	return b
}

// Build constructs an immutable Matcher from the patterns added to the builder.
// The resulting Matcher is safe for concurrent use without locks.
//
// Returns an error if:
//   - no patterns were added
//   - all patterns are empty
func (b *Builder) Build() (*Matcher, error) {
	if len(b.patterns) == 0 {
		return nil, errors.New("no patterns provided")
	}
	
	m := &Matcher{
		root: &Node{
			depth: 0,
		},
		patterns:    make([]string, 0, len(b.patterns)),
		patternLens: make([]int, 0, len(b.patterns)),
	}
	
	for _, pattern := range b.patterns {
		m.addPattern(pattern, len(m.patterns))
		m.patterns = append(m.patterns, pattern)
		m.patternLens = append(m.patternLens, utf8.RuneCountInString(pattern))
	}
	
	m.buildFailureFunction()
	return m, nil
}

// Node represents a node in the Aho-Corasick trie data structure.
// Each node contains references to its children, parent, and failure link,
// as well as output patterns that end at this node.
type Node struct {
	// asciiChildren provides fast access for ASCII characters (0-127)
	asciiChildren [128]*Node
	// children stores non-ASCII characters using a map
	children      map[rune]*Node
	// parent points to the parent node in the trie
	parent        *Node
	// fail points to the failure link used during matching
	fail          *Node
	// output contains indices of patterns that end at this node
	output        []int
	// depth is the distance from the root node
	depth         int
	// character is the character that leads to this node from its parent
	character     rune
}

// getChild returns the child node for the given rune, or nil if not found.
// Uses optimized array lookup for ASCII characters.
func (n *Node) getChild(r rune) *Node {
	if r < 128 {
		return n.asciiChildren[r]
	}
	return n.children[r]
}

// setChild sets the child node for the given rune.
// Uses optimized array storage for ASCII characters and map for others.
func (n *Node) setChild(r rune, child *Node) {
	if r < 128 {
		n.asciiChildren[r] = child
	} else {
		if n.children == nil {
			n.children = make(map[rune]*Node)
		}
		n.children[r] = child
	}
}

// forEachChild iterates over all child nodes and calls the provided function.
// Iterates through ASCII children first, then non-ASCII children.
func (n *Node) forEachChild(fn func(rune, *Node)) {
	for i, child := range n.asciiChildren {
		if child != nil {
			fn(rune(i), child)
		}
	}
	for r, child := range n.children {
		fn(r, child)
	}
}

// Matcher implements the Aho-Corasick algorithm for multiple pattern matching.
// It builds a trie from patterns and uses failure links for efficient matching.
// Once built, a Matcher is immutable and safe for concurrent use without locks.
type Matcher struct {
	// root is the root node of the trie
	root        *Node
	// patterns stores all patterns indexed by their ID
	patterns    []string
	// patternLens stores the rune length of each pattern for efficient access
	patternLens []int
}

// Match represents a pattern match found in the text.
type Match struct {
	// Pattern is the matched pattern string
	Pattern string
	// Index is the index of the pattern in the original pattern slice
	Index   int
	// Start is the starting position of the match in the text (in runes)
	Start   int
	// End is the ending position of the match in the text (in runes)
	End     int
}

// FindAll finds all occurrences of the patterns in the given text.
// Returns a slice of Match structs, each representing a pattern match.
// The matches are returned in the order they are found in the text.
// Overlapping matches are all reported.
//
// The Start and End positions in Match are measured in runes, not bytes.
// This method properly handles UTF-8 encoded text.
// This method is safe for concurrent use.
//
// Returns an error if:
//   - text contains invalid UTF-8 sequences
func (m *Matcher) FindAll(text string) ([]Match, error) {
	if m == nil || m.root == nil {
		return nil, errors.New("matcher not built")
	}
	
	if !utf8.ValidString(text) {
		return nil, errors.New("text contains invalid UTF-8 sequences")
	}
	
	matches := make([]Match, 0, 16)
	node := m.root
	
	pos := 0
	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])
		
		for node != m.root && node.getChild(r) == nil {
			node = node.fail
		}
		
		if child := node.getChild(r); child != nil {
			node = child
		}
		
		if len(node.output) > 0 {
			for _, patternIndex := range node.output {
				pattern := m.patterns[patternIndex]
				patternRuneLen := m.patternLens[patternIndex]
				matches = append(matches, Match{
					Pattern: pattern,
					Index:   patternIndex,
					Start:   pos - patternRuneLen + 1,
					End:     pos + 1,
				})
			}
		}
		
		i += size
		pos++
	}
	
	return matches, nil
}

// addPattern adds a pattern to the trie with the given index.
// It creates nodes as needed and marks the final node with the pattern index.
func (m *Matcher) addPattern(pattern string, index int) {
	node := m.root
	
	for _, ch := range pattern {
		child := node.getChild(ch)
		if child == nil {
			child = &Node{
				parent:    node,
				depth:     node.depth + 1,
				character: ch,
			}
			node.setChild(ch, child)
		}
		node = child
	}
	
	if node.output == nil {
		node.output = make([]int, 0)
	}
	node.output = append(node.output, index)
}

// buildFailureFunction constructs failure links for the Aho-Corasick automaton.
// This enables the algorithm to efficiently handle mismatches by falling back
// to the longest proper suffix that is also a prefix of some pattern.
func (m *Matcher) buildFailureFunction() {
	queue := make([]*Node, 0, 64)
	
	m.root.forEachChild(func(ch rune, child *Node) {
		child.fail = m.root
		queue = append(queue, child)
	})
	
	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		
		currentNode.forEachChild(func(ch rune, child *Node) {
			queue = append(queue, child)
			
			failNode := currentNode.fail
			for failNode != nil && failNode.getChild(ch) == nil {
				failNode = failNode.fail
			}
			
			if failNode == nil {
				child.fail = m.root
			} else {
				child.fail = failNode.getChild(ch)
				if child.fail == child {
					child.fail = m.root
				}
			}
			
			if child.fail != nil && len(child.fail.output) > 0 {
				child.output = append(child.output, child.fail.output...)
			}
		})
	}
}