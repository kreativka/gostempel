package gostempel // import "github.com/kreativka/gostempel"

import (
	"os"
	"strings"

	"github.com/kreativka/gostempel/egothor"
	"github.com/kreativka/gostempel/javautf"
)

// Trie interface for egothor trie
type Trie interface {
	GetLastOnPath([]rune) []rune
}

// Minimum length of token
const minTokenLength = 3

// Stem returns stem from given token
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

	// Read method from stem file
	method, err := javautf.ReadUTF(f)
	if err != nil {
		return nil, err
	}

	var rv Trie
	if strings.HasPrefix(method, "-0ME2") {
		rv, err = egothor.NewMultiTrie(f)
		if err != nil {
			return nil, err
		}
	} else {
		rv, err = egothor.NewTrie(f)
		if err != nil {
			return nil, err
		}
	}
	return rv, nil
}
