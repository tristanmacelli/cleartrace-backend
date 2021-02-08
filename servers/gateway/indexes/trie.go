package indexes

import (
	"sync"
)

//TODO: implement a trie data structure that stores
//keys of type string and values of type int64

//PRO TIP: if you are having troubles and want to see
//what your trie structure looks like at various points,
//either use the debugger, or try this package:
//https://github.com/davecgh/go-spew

type trieNode struct {
	children map[rune]*trieNode
	values   int64set
}

//Trie implements a trie data structure mapping strings to int64s
//that is safe for concurrent use.
type Trie struct {
	Root *trieNode
	lock *sync.Mutex
}

//NewTrie constructs a new Trie.
func NewTrie(newLock *sync.Mutex) *Trie {
	return &Trie{
		Root: &trieNode{},
		lock: newLock,
	}
}

//Len returns the number of entries in the trie.
func (t *Trie) Len() int {
	return lenHelper(t.Root)
}

func lenHelper(node *trieNode) int {
	if len(node.children) == 0 {
		return len(node.values.all())
	} else {
		sum := len(node.values.all())
		for _, child := range node.children {
			sum += lenHelper(child)
		}
		return sum
	}
}

//Add adds a key and value to the trie.
func (t *Trie) Add(key string, value int64) {
	root := t.Root
	for _, r := range key {
		child, _ := root.children[r]
		if child == nil {
			if root.children == nil {
				root.children = map[rune]*trieNode{}
			}
			child = &trieNode{}
			root.children[r] = child
		}
		root = child
	}
	root.values.add(value)
}

//Find finds `max` values matching `prefix`. If the trie
//is entirely empty, or the prefix is empty, or max == 0,
//or the prefix is not found, this returns a nil slice.
func (t *Trie) Find(prefix string, max int) []int64 {
	if len(prefix) == 0 || max == 0 || len(t.Root.children) == 0 {
		return []int64{}
	}
	var ids []int64
	node := t.Root
	// Loop to the end of prefix, returning a nil slice if the prefix isn't present
	for _, r := range prefix {
		child, _ := node.children[r]
		if child == nil {
			return []int64{}
		}
		node = child
	}
	return findHelper(node, max, ids)
}

func findHelper(node *trieNode, max int, ids []int64) []int64 {
	if len(node.values.all()) == 0 && len(node.children) == 0 {
		return ids
	} else {
		for _, id := range node.values.all() {
			ids = append(ids, id)
			if len(ids) >= max {
				return ids
			}
		}
		for _, child := range node.children {
			return findHelper(child, max, ids)
		}
		return ids
	}
}

//Remove removes a key/value pair from the trie
//and trims branches with no values.
func (t *Trie) Remove(key string, value int64) {
	node := t.Root
	for _, r := range key {
		child, _ := node.children[r]
		if child == nil {
			// Since the key does not exist there is nothing to delete
			return
		}
		node = child
	}
	node.values.remove(value)
	ids := node.values.all()
	// If the node contains other values or has children the tree doesn't need trimming
	if len(ids) > 0 || !node.isLeaf() {
		return
	}
	node = t.Root
	t.Root = removeHelper(node, []rune(key), len(key))
}

func removeHelper(node *trieNode, runes []rune, length int) *trieNode {
	if len(runes) == 0 {
		return nil
	} else {
		child, _ := node.children[runes[0]]
		newRunes := []rune{}
		if len(runes) > 1 {
			newRunes = runes[1:]
		}
		node.children[runes[0]] = removeHelper(child, newRunes, length)
		if node.children[runes[0]] == nil && len(node.children) == 1 {
			node.children = nil
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
