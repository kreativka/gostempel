package egothor

import (
	"fmt"

	"github.com/kreativka/reader/stringutil"
)

// Trie s
type Trie struct {
	Eom     string
	forward bool
	cmds    []string
	rows    []*Row
	root    int32
}

// NewTrie returns trie
func NewTrie() *Trie {
	return &Trie{}
}

// AddRow adds row
func (t *Trie) AddRow(row *Row) {
	t.rows = append(t.rows, row)
}

// AddCmds adds cmd to cmds
func (t *Trie) AddCmds(cmd string) {
	t.cmds = append(t.cmds, cmd)
}

// Cmds returns cmds
func (t Trie) Cmds() []string {
	return t.cmds
}

// SetForward sets forward
func (t *Trie) SetForward(forward bool) {
	t.forward = forward
}

// Forward returns forward
func (t Trie) Forward() bool {
	return t.forward
}

// SetRoot sets root
func (t *Trie) SetRoot(root int32) {
	t.root = root
}

// GetLastOnPath returns something
func (t Trie) GetLastOnPath(key string) ([]rune, bool) {
	var ok bool
	var last []rune
	now, err := t.getRow(t.root)
	if err != nil {
		return last, ok
	}

	var w int32
	if !t.Forward() {
		key = stringutil.Reverse(key)
	}
	for i, c := range key {
		w = now.Cmd(c)
		if w >= 0 {
			last = []rune(t.cmds[w])
			ok = true
		}
		// Return from end of original function
		// w = now.getCmd(new Character(e.next()));
		// return (w >= 0) ? cmds.elementAt(w) : last;
		if i > len(key) {
			return last, ok
		}
		w = now.Ref(c)
		if w >= 0 {
			now, err = t.getRow(w)
			if err != nil {
				return last, ok
			}
		} else {
			return last, ok
		}
	}
	return last, ok
}

// getRow returns row
func (t Trie) getRow(i int32) (*Row, error) {
	var err error
	if i < 0 || i >= int32(len(t.rows)) {
		err = fmt.Errorf("row %d not found not found", i)
		return nil, err
	}
	return t.rows[i], err
}
