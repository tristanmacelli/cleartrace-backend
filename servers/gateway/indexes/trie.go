package indexes

import (
	"fmt"
	"strings"
	"sync"
)

// TODO: implement a trie data structure that stores
// keys of type string and values of type int64

// PRO TIP: if you are having troubles and want to see
// what your trie structure looks like at various points,
// either use the debugger, or try this package:
// https://github.com/davecgh/go-spew

type trieNode struct {
	children map[rune]*trieNode
	values   int64set
}

// Trie implements a trie data structure mapping strings to int64s
// that is safe for concurrent use.
type Trie struct {
	Root *trieNode
	lock *sync.Mutex
}

// NewTrie constructs a new Trie.
func NewTrie(newLock *sync.Mutex) *Trie {
	return &Trie{
		Root: &trieNode{children: map[rune]*trieNode{}, values: int64set{}},
		lock: newLock,
	}
}

// Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	node := t.Root
	return lenHelper(node)
}

func lenHelper(node *trieNode) int {
	if node.isLeaf() {
		return len(node.values.all())
	} else {
		sum := len(node.values.all())
		for _, child := range node.children {
			if child != nil {
				sum += lenHelper(child)
			}
		}
		return sum
	}
}

// Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) error {
	if key == "" {
		return fmt.Errorf("error: must enter a non-empty string")
	}
	if value < 0 {
		return fmt.Errorf("error: invalid id value. Must be non-negative")
	}
	key = strings.ToLower(key)
	keys := strings.Fields(key)
	// If a key contains spaces it is separated into multiple keys
	for _, key := range keys {
		node := t.Root
		for _, r := range key {
			// Check that this still works without (, _)
			child := node.children[r]
			if child == nil {
				child = &trieNode{children: map[rune]*trieNode{}, values: int64set{}}
				t.lock.Lock()
				node.children[r] = child
				t.lock.Unlock()
			}
			node = child
		}
		t.lock.Lock()
		node.values.add(value)
		t.lock.Unlock()
	}
	return nil
}

// Find finds `max` values matching `prefix`. If the trie
// is entirely empty, or the prefix is empty, or max == 0,
// or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) ([]int64, error) {
	if max < 0 {
		return nil, fmt.Errorf("error: invalid max value. Must be non-negative")
	}
	prefix = strings.ToLower(prefix)
	if len(prefix) == 0 || max == 0 || len(t.Root.children) == 0 {
		return []int64{}, nil
	}
	ids := make(int64set)
	for i, prefixPart := range strings.Fields(prefix) {
		node := t.Root
		// Loop to the end of prefix, returning a nil slice if the prefix isn't present
		for _, r := range prefixPart {
			// Check that this still works without (, _)
			child := node.children[r]
			if child == nil {
				return []int64{}, nil
			}
			node = child
		}
		newIds := make(int64set)
		subsequent := i > 0
		ids = findHelper(node, max, ids, newIds, subsequent)
	}
	return ids.all(), nil
}

func findHelper(node *trieNode, max int, ids int64set, newIds int64set, subsequentToken bool) int64set {
	if len(node.values.all()) == 0 && len(node.children) == 0 {
		return ids
	} else {
		for _, id := range node.values.all() {
			// Adds to union list if method is called on a subsequent token
			// otherwise adds directly to ids slice until the max is reached
			if subsequentToken && ids.has(id) {
				newIds.add(id)
			} else if !subsequentToken {
				ids.add(id)
			}
			if len(newIds.all()) >= max {
				return newIds
			}
			if len(ids.all()) >= max {
				return ids
			}
		}
		for _, child := range node.children {
			if child != nil {
				return findHelper(child, max, ids, newIds, subsequentToken)
			}
		}
		return ids
	}
}

// Remove removes a key/value pair from the trie
// and trims branches with no values.
func (t *Trie) Remove(key string, value int64) error {
	if key == "" {
		return fmt.Errorf("error: must enter a non-empty string")
	}
	if value < 0 {
		return fmt.Errorf("error: invalid id value")
	}
	key = strings.ToLower(key)
	node := t.Root
	for _, r := range key {
		// Check that this still works without (, _)
		child := node.children[r]
		if child == nil {
			// Since the key does not exist there is nothing to delete
			return nil
		}
		node = child
	}
	t.lock.Lock()
	node.values.remove(value)
	t.lock.Unlock()
	ids := node.values.all()
	// If the node contains other values or has children the tree doesn't need trimming
	if len(ids) > 0 || !node.isLeaf() {
		return nil
	}
	node = t.Root
	t.Root = t.removeHelper(node, []rune(key), len(key))
	return nil
}

func (t *Trie) removeHelper(node *trieNode, runes []rune, length int) *trieNode {
	if len(runes) == 0 {
		return nil
	} else {
		// Check that this still works without (, _)
		child := node.children[runes[0]]
		newRunes := []rune{}
		if len(runes) > 1 {
			newRunes = runes[1:]
		}
		child = t.removeHelper(child, newRunes, length)
		t.lock.Lock()
		node.children[runes[0]] = child
		t.lock.Unlock()
		if node.children[runes[0]] == nil && len(node.children) == 1 {
			t.lock.Lock()
			node.children = nil
			t.lock.Unlock()
		}
		if len(runes) < length && node.isLeaf() && len(node.values) == 0 {
			return nil
		}
		return node
	}
}

func (t *trieNode) isLeaf() bool {
	return len(t.children) == 0
}
