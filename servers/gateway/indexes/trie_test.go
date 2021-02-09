package indexes

import (
	"sync"
	"testing"
)

func callLen(t *Trie, additions bool, removals bool) int {
	if additions {
		t.Add("First", 1)
		t.Add("Second", 2)
	}
	if removals {
		t.Remove("First", 1)
	}
	return t.Len()
}

//TODO: implement automated tests for your trie data structure
func TestLen(t *testing.T) {
	cases := []struct {
		name      string
		hint      string
		t         *Trie
		additions bool
		removals  bool
		length    int
		expected  string
	}{
		{
			"Checking empty trie returns 0 length when calling len",
			"Make sure you are not counting the root node",
			NewTrie(&sync.Mutex{}),
			false,
			false,
			0,
			"A length greater than 0 was not expected but was returned",
		},
		{
			"Checking trie with additions returns accurate lenght",
			"Make sure you are counting entries/additions and not depth",
			NewTrie(&sync.Mutex{}),
			true,
			false,
			0,
			"A length other than 2 was not expected but was returned",
		},
		{
			"Checking trie with removals returns accurate length",
			"Make sure you are removing the value and the nodes",
			NewTrie(&sync.Mutex{}),
			true,
			true,
			0,
			"A length greater than 0 was not expected but was returned",
		},
	}
	for _, c := range cases {
		rr := callLen(c.t, c.additions, c.removals)
		if c.length != rr {
			t.Errorf(c.expected)
		}
	}
}

func TestAdd(t *testing.T) {}

func TestFind(t *testing.T) {}

func TestRemove(t *testing.T) {}
