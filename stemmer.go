package gostempel // import "github.com/kreativka/gostempel"

import (
	"fmt"
	"os"
	"strings"

	"github.com/blevesearch/stempel/javadata"
	"github.com/kreativka/gostempel/egothor"
)

// Trie interface for egothor trie
type Trie interface {
	GetLastOnPath([]rune) []rune
}

// Minimum length of token
const minTokenLength = 3

// Stem returns stem from a given token
func Stem(stem Trie, token []rune) []rune {
	// Don't stem tokens less than 3 chars and empty tokens
	if token == nil || len(token) <= minTokenLength {
		return token
	}

	// Get commands to create stem
	cmd := stem.GetLastOnPath(token)
	if cmd == nil {
		return token
	}

	// Apply cmds to token and return stem
	res := egothor.DiffApply(token, cmd)
	if len(res) > 0 {
		return res
	}
	return token
}

// LoadStemmer returns MultiTrie from given file
func LoadStemmer(filename string) (Trie, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	r := javadata.NewReader(f)
	method, err := r.ReadUTF()
	if err != nil {
		return nil, fmt.Errorf("error reading method value: %v", err)
	}

	var rv Trie
	if strings.HasPrefix(method, "-0ME2") {
		rv, err = egothor.NewMultiTrie(r)
		if err != nil {
			return nil, fmt.Errorf("error creating egothor trie: %v", err)
		}
	} else {
		rv, err = egothor.NewTrie(r)
		if err != nil {
			return nil, fmt.Errorf("error creating egothor trie: %v", err)
		}
	}
	return rv, nil
}
