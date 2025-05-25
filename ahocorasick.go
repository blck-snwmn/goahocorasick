package goahocorasick

import (
	"container/list"
)

type Node struct {
	children  map[rune]*Node
	parent    *Node
	fail      *Node
	output    []int
	depth     int
	character rune
}

type Matcher struct {
	root     *Node
	patterns []string
}

func New() *Matcher {
	return &Matcher{
		root: &Node{
			children: make(map[rune]*Node),
			depth:    0,
		},
		patterns: make([]string, 0),
	}
}

func (m *Matcher) Build(patterns []string) {
	m.root = &Node{
		children: make(map[rune]*Node),
		depth:    0,
	}
	m.patterns = make([]string, 0)
	
	for _, pattern := range patterns {
		if len(pattern) > 0 {
			m.addPattern(pattern, len(m.patterns))
			m.patterns = append(m.patterns, pattern)
		}
	}
	
	m.buildFailureFunction()
}

func (m *Matcher) addPattern(pattern string, index int) {
	node := m.root
	
	for _, ch := range pattern {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = &Node{
				children:  make(map[rune]*Node),
				parent:    node,
				depth:     node.depth + 1,
				character: ch,
			}
		}
		node = node.children[ch]
	}
	
	if node.output == nil {
		node.output = make([]int, 0)
	}
	node.output = append(node.output, index)
}

func (m *Matcher) buildFailureFunction() {
	queue := list.New()
	
	for _, child := range m.root.children {
		child.fail = m.root
		queue.PushBack(child)
	}
	
	for queue.Len() > 0 {
		element := queue.Front()
		currentNode := element.Value.(*Node)
		queue.Remove(element)
		
		for ch, child := range currentNode.children {
			queue.PushBack(child)
			
			failNode := currentNode.fail
			for failNode != nil && failNode.children[ch] == nil {
				failNode = failNode.fail
			}
			
			if failNode == nil {
				child.fail = m.root
			} else {
				child.fail = failNode.children[ch]
				if child.fail == child {
					child.fail = m.root
				}
			}
			
			if child.fail != nil && len(child.fail.output) > 0 {
				child.output = append(child.output, child.fail.output...)
			}
		}
	}
}

type Match struct {
	Pattern string
	Index   int
	Start   int
	End     int
}

func (m *Matcher) FindAll(text string) []Match {
	matches := make([]Match, 0)
	node := m.root
	
	runes := []rune(text)
	for i, ch := range runes {
		for node != m.root && node.children[ch] == nil {
			node = node.fail
		}
		
		if node.children[ch] != nil {
			node = node.children[ch]
		}
		
		if len(node.output) > 0 {
			for _, patternIndex := range node.output {
				pattern := m.patterns[patternIndex]
				matches = append(matches, Match{
					Pattern: pattern,
					Index:   patternIndex,
					Start:   i - len([]rune(pattern)) + 1,
					End:     i + 1,
				})
			}
		}
	}
	
	return matches
}