package goahocorasick

import (
	"unicode/utf8"
)

type Node struct {
	asciiChildren [128]*Node
	children      map[rune]*Node
	parent        *Node
	fail          *Node
	output        []int
	depth         int
	character     rune
}

func (n *Node) getChild(r rune) *Node {
	if r < 128 {
		return n.asciiChildren[r]
	}
	return n.children[r]
}

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

type Matcher struct {
	root        *Node
	patterns    []string
	patternLens []int
}

func New() *Matcher {
	return &Matcher{
		root: &Node{
			depth: 0,
		},
		patterns:    make([]string, 0),
		patternLens: make([]int, 0),
	}
}

func (m *Matcher) Build(patterns []string) {
	m.root = &Node{
		depth: 0,
	}
	m.patterns = make([]string, 0)
	m.patternLens = make([]int, 0)
	
	for _, pattern := range patterns {
		if len(pattern) > 0 {
			m.addPattern(pattern, len(m.patterns))
			m.patterns = append(m.patterns, pattern)
			m.patternLens = append(m.patternLens, utf8.RuneCountInString(pattern))
		}
	}
	
	m.buildFailureFunction()
}

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

type Match struct {
	Pattern string
	Index   int
	Start   int
	End     int
}

func (m *Matcher) FindAll(text string) []Match {
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
	
	return matches
}