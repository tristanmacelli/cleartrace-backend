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

func callAdd(t *Trie, additions []string) int {
	for id, str := range additions {
		t.Add(str, int64(id))
	}
	return t.Len()
}

func callRemove(t *Trie, additions []string, removals []string, ids []int64) int {
	for id, str := range additions {
		t.Add(str, int64(id))
	}
	for id, str := range removals {
		t.Add(str, ids[id])
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
			"Checking trie with additions returns accurate length",
			"Make sure you are counting entries/additions and not depth",
			NewTrie(&sync.Mutex{}),
			true,
			false,
			2,
			"A length other than 2 was not expected but was returned",
		},
		{
			"Checking trie with removals returns accurate length",
			"Make sure you are removing the value and the nodes",
			NewTrie(&sync.Mutex{}),
			true,
			true,
			1,
			"A length greater than 0 was not expected but was returned",
		},
	}
	for _, c := range cases {
		rr := callLen(c.t, c.additions, c.removals)
		if c.length != rr {
			t.Errorf("case %s: unexpected error %v\nHINT: %s", c.name, c.expected, c.hint)
		}
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name      string
		hint      string
		t         *Trie
		additions []string
		length    int
		expected  string
	}{
		{
			"Check if adding to an empty tree works",
			"The length should be greater than 0",
			NewTrie(&sync.Mutex{}),
			[]string{"Bob"},
			1,
			"A length of 0 was not expected but was returned",
		},
		{
			"Check if adding to a unique key to the tree works",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Muna", "Abby"},
			2,
			"A length of 1 was not expected but was returned",
		},
		{
			"Check if adding to a non-unique key to the tree works (full-prefix match)",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Ngugyen", "Ngugyen"},
			2,
			"A length of 1 was not expected but was returned",
		},
		{
			"Check if adding to a non-unique key to the tree works (partial prefix match)",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Jon", "Jonathan"},
			2,
			"A length of 1 was not expected but was returned",
		},
		{
			"Check if adding to a non-unique key to  tree works (partial prefix match)",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Jonathan", "Jon"},
			2,
			"A length of 1 was not expected but was returned",
		},
	}
	for _, c := range cases {
		rr := callAdd(c.t, c.additions)
		if c.length != rr {
			t.Errorf(c.expected)
		}
	}
}

func TestFind(t *testing.T) {}

func TestRemove(t *testing.T) {
	cases := []struct {
		name      string
		hint      string
		t         *Trie
		additions []string
		removals  []string
		ids       []int64
		length    int
		expected  string
	}{
		{
			"Check if removing last item in the tree works",
			"The length should be less than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"George"},
			[]string{"George"},
			[]int64{1},
			1,
			"A length greater than 0 was not expected but was returned",
		},
		{
			"Check if removing multiple items from the tree works",
			"The length should be less than 2",
			NewTrie(&sync.Mutex{}),
			[]string{"Max", "Allison", "Cooper"},
			[]string{"Max", "Cooper"},
			[]int64{1, 3},
			1,
			"A length greater than 1 was not expected but was returned",
		},
		{
			"Check if removing a unique item in the tree works",
			"The length should be less than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Seaeun", "Mitchell"},
			[]string{"Seaeun"},
			[]int64{1},
			1,
			"A length greater than 1 was not expected but was returned",
		},
		{
			"Check if removing a non-unique item from the tree works (full key match)",
			"The length should be less than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Connor", "Connor"},
			[]string{"Connor"},
			[]int64{1},
			1,
			"A length greater than 1 was not expected but was returned",
		},
		{
			"Check if removing a non-unique item from the tree works (partial key match)",
			"The length should be less than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Melissa", "Melinda"},
			[]string{"Melissa"},
			[]int64{1},
			1,
			"A length greater than 1 was not expected but was returned",
		},
		{
			"Check if removing a non-unique key from the tree works (partial key match)",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Jon", "Jonathan"},
			[]string{"Jon"},
			[]int64{1},
			2,
			"A length of 1 was not expected but was returned",
		},
		{
			"Check if removing a non-unique key from the tree works (partial key match)",
			"The length should be greater than 1",
			NewTrie(&sync.Mutex{}),
			[]string{"Jon", "Jonathan"},
			[]string{"Jonathan"},
			[]int64{2},
			2,
			"A length of 1 was not expected but was returned",
		},
	}
	for _, c := range cases {
		rr := callRemove(c.t, c.additions, c.removals, c.ids)
		if c.length != rr {
			t.Errorf(c.expected)
		}
	}
}
