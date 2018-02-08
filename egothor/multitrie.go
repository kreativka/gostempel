package egothor

import (
	"strings"
)

// MultiTrie struct
type MultiTrie struct {
	EOM     rune
	forward bool
	BY      int32
	cmds    []string
	root    int32
	tries   []*Trie
}

// AddCmds adds cmd to cmds
func (t *MultiTrie) AddCmds(cmd string) {
	t.cmds = append(t.cmds, cmd)
}

// Cmds returns cmds
func (t MultiTrie) Cmds() []string {
	return t.cmds
}

// AddTrie adds
func (t *MultiTrie) AddTrie(trie *Trie) {
	t.tries = append(t.tries, trie)
}

// SetForward sets forward
func (t *MultiTrie) SetForward(forward bool) {
	t.forward = forward
}

// SetBY sets by :)
func (t *MultiTrie) SetBY(by int32) {
	t.BY = by
}

// SetRoot sets root
func (t *MultiTrie) SetRoot(root int32) {
	t.root = root
}

// NewMultiTrie return Trie
func NewMultiTrie() *MultiTrie {
	return &MultiTrie{EOM: '*', forward: false, root: 0, BY: 0}
}

func cannotFollow(after, goes rune) bool {
	switch after {
	case '-':
		return after == goes
	case 'D':
		return after == goes
	}
	return false
}

func (t MultiTrie) skip(in string, count int) (string, bool) {
	runes := []rune(in)

	if len(runes)-count < 0 {
		return "", false
	}
	if t.forward {
		return string(runes[count:]), true
	}
	return string(runes[:len(runes)-count]), true
}

func lengthPP(cmd string) int {
	runes := []rune(cmd)
	l := 0
	for i, c := range runes {
		i++
		switch c {
		case '-':
			l += int(runes[i]) - 'a' + 1
			break
		case 'D':
			l += int(runes[i]) - 'a' + 1
			break
		case 'R':
			l++
			break
		case 'I':
			break
			// }
		}
	}
	return l
}

// GetLastOnPath returns something
func (t *MultiTrie) GetLastOnPath(key string) string {
	var result string

	lastkey := key
	lastch := rune(' ')
	p := make(map[int]string)

	for i := 0; i < len(t.tries); i++ {
		r, ok := t.tries[i].GetLastOnPath(lastkey)
		if !ok || (r[0] == t.EOM && len(r) == 1) {
			return result
		}

		if cannotFollow(lastch, r[0]) {
			return result
		}

		lastch = r[len(r)-2]

		p[i] = string(r)
		if strings.HasPrefix(p[i], "-") {
			var ok bool
			if i > 0 {
				key, ok = t.skip(key, lengthPP(p[i-1]))
				if !ok {
					return result
				}
			}
			key, ok = t.skip(key, lengthPP(p[i]))
			if !ok {
				return result
			}
		}
		result += string(r)
		if len(key) != 0 {
			lastkey = key
		}
	}
	return result
}
