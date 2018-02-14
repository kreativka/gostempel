package gostempel

import (
	"bufio"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/kreativka/gostempel/egothor"
	"github.com/kreativka/gostempel/javautf"
)

// Minimum length of token
const minTokenLength = 3

// Stem returns stem from given token
func Stem(stem egothor.Tries, token string) string {
	// Don't stem tokens less than 3 chars and empty tokens
	if token == "" || utf8.RuneCountInString(token) <= minTokenLength {
		return token
	}

	// Get commands to create stem
	cmd, ok := stem.GetLastOnPath(token)
	if cmd == nil || !ok {
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
func LoadStemmer(filename string) (egothor.Tries, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	in := bufio.NewReader(f)

	// Read method from stem file
	m, err := javautf.ReadUTF(in)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(m, "-0ME2") {
		return egothor.NewMultiTrie(in), nil
	}
	return egothor.NewTrie(in), nil
}
