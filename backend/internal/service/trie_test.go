package service

import (
	"slices"
	"testing"
)

func newTestTrie(names ...string) *trie {
	tr := newTrie()
	for _, name := range names {
		tr.insert(name)
	}
	return tr
}

func TestTrieSearch_PrefixMatch(t *testing.T) {
	tr := newTestTrie(
		"Red Dead Redemption",
		"Red Dead Redemption 2",
		"Reddit Simulator",
		"The Witcher 3",
	)

	got := tr.search("Red ", 20)
	want := []string{"Red Dead Redemption", "Red Dead Redemption 2"}
	if !slices.Equal(got, want) {
		t.Errorf("search(%q, 20) = %v, want %v", "Red ", got, want)
	}
}

func TestTrieSearch_CaseInsensitive(t *testing.T) {
	tr := newTestTrie("Red Dead Redemption")

	got := tr.search("red ", 20)
	want := []string{"Red Dead Redemption"}
	if !slices.Equal(got, want) {
		t.Errorf("search(%q, 20) = %v, want %v", "red ", got, want)
	}
}

func TestTrieSearch_NoMatchReturnsNil(t *testing.T) {
	tr := newTestTrie("Red Dead Redemption")

	if got := tr.search("halo", 20); got != nil {
		t.Errorf("search(%q, 20) = %v, want nil", "halo", got)
	}
}

func TestTrieSearch_RespectsLimit(t *testing.T) {
	tr := newTestTrie("Aa", "Ab", "Ac", "Ad", "Ae")

	got := tr.search("A", 3)
	want := []string{"Aa", "Ab", "Ac"}
	if !slices.Equal(got, want) {
		t.Errorf("search(%q, 3) = %v, want %v", "A", got, want)
	}
}

func TestTrieSearch_AlphabeticalOrder(t *testing.T) {
	tr := newTestTrie("Ghost Zelda", "Ghost Amnesia", "Ghost Minecraft")

	got := tr.search("Ghost", 20)
	want := []string{"Ghost Amnesia", "Ghost Minecraft", "Ghost Zelda"}
	if !slices.Equal(got, want) {
		t.Errorf("search(%q, 20) = %v, want %v", "Ghost", got, want)
	}
}
