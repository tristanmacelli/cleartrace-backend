package indexes

import "testing"

//TODO: implement automated tests for your trie data structure
func TestLen(t *testing.T) {
	cases := []struct {
		name     string
		hint     string
		length   int
		expected string
	}{
		{
			"Checking empty trie returns 0 length when calling len",
			"Make sure you are not counting the root node",
			0,
			"A length greater than 0 was not expected but was returned",
		},
	}
	for _, c := range cases {
		// rr := len()
		if c.length != 0 {
			t.Errorf(c.expected)
		}
	}
}

func TestAdd(t *testing.T) {}

func TestFind(t *testing.T) {}

func TestRemove(t *testing.T) {}
