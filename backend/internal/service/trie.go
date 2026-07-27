package service

import (
	"sort"
	"strings"
)

type trieNode struct {
	children map[rune]*trieNode
	// matches holds every name that passes through this node, i.e. every
	// name whose prefix (up to this depth) equals the path from the root.
	// Populated at insert time so a lookup is just a walk to the node.
	matches []string
}

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[rune]*trieNode)}
}

type trie struct {
	root *trieNode
}

func newTrie() *trie {
	return &trie{root: newTrieNode()}
}

func (t *trie) insert(name string) {
	node := t.root
	for _, r := range strings.ToLower(name) {
		child, ok := node.children[r]
		if !ok {
			child = newTrieNode()
			node.children[r] = child
		}
		node = child
		node.matches = append(node.matches, name)
	}
}

// search returns up to limit names starting with prefix (case-insensitive),
// sorted alphabetically. It returns nil if no name matches.
func (t *trie) search(prefix string, limit int) []string {
	node := t.root
	for _, r := range strings.ToLower(prefix) {
		child, ok := node.children[r]
		if !ok {
			return nil
		}
		node = child
	}

	results := append([]string(nil), node.matches...)
	sort.Strings(results)
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}
