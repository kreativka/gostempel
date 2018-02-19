// Egothor Software License version 1.00
// Copyright (C) 1997-2004 Leo Galambos.
// Copyright (C) 2002-2004 "Egothor developers"
//   on behalf of the Egothor Project.
//          All rights reserved.
// This  software  is  copyrighted  by  the "Egothor developers". If this
// license applies to a single file or document, the "Egothor developers"
// are the people or entities mentioned as copyright holders in that file
// or  document.  If  this  license  applies  to the Egothor project as a
// whole,  the  copyright holders are the people or entities mentioned in
// the  file CREDITS. This file can be found in the same location as this
// license in the distribution.
// Redistribution  and  use  in  source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
// 1. Redistributions  of  source  code  must retain the above copyright
// notice, the list of contributors, this list of conditions, and the
// following disclaimer.
// 2. Redistributions  in binary form must reproduce the above copyright
// notice, the list of contributors, this list of conditions, and the
// disclaimer  that  follows  these  conditions  in the documentation
// and/or other materials provided with the distribution.
// 3. The name "Egothor" must not be used to endorse or promote products
// derived  from  this software without prior written permission. For
// written permission, please contact Leo.G@seznam.cz
// 4. Products  derived  from this software may not be called "Egothor",
// nor  may  "Egothor"  appear  in  their name, without prior written
// permission from Leo.G@seznam.cz.
// In addition, we request that you include in the end-user documentation
// provided  with  the  redistribution  and/or  in the software itself an
// acknowledgement equivalent to the following:
// "This product includes software developed by the Egothor Project.
// http://egothor.sf.net/"
// THIS  SOFTWARE  IS  PROVIDED  ``AS  IS''  AND ANY EXPRESSED OR IMPLIED
// WARRANTIES,  INCLUDING,  BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
// MERCHANTABILITY  AND  FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN  NO  EVENT  SHALL THE EGOTHOR PROJECT OR ITS CONTRIBUTORS BE LIABLE
// FOR   ANY   DIRECT,   INDIRECT,  INCIDENTAL,  SPECIAL,  EXEMPLARY,  OR
// CONSEQUENTIAL  DAMAGES  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE  GOODS  OR  SERVICES;  LOSS  OF  USE,  DATA, OR PROFITS; OR
// BUSINESS  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER  IN  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
// IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// This  software  consists  of  voluntary  contributions  made  by  many
// individuals  on  behalf  of  the  Egothor  Project  and was originally
// created by Leo Galambos (Leo.G@seznam.cz).

package egothor

import (
	"fmt"
	"io"

	"github.com/kreativka/gostempel/javautf"
)

// Tries wraps tries
type Tries interface {
	GetLastOnPath([]rune) ([]rune, bool)
}

// Trie implements Tries interface
type Trie struct {
	cmds    [][]rune
	forward bool
	root    int32
	rows    []*Row
}

// NewTrie returns trie
func NewTrie(r io.Reader) (*Trie, error) {
	t := Trie{}
	br := &errBinaryReader{r: r}
	br.Read(&t.forward)
	br.Read(&t.root)

	var i int32
	br.Read(&i)
	for ; i > 0; i-- {
		cmd, err := javautf.ReadUTF(r)
		if err != nil {
			return nil, err
		}
		t.AddCmds([]rune(cmd))
	}

	br.Read(&i)
	for ; i > 0; i-- {
		nr, err := NewRow(r)
		if err != nil {
			return nil, err
		}

		t.AddRow(nr)
	}
	err := br.Err()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// AddRow adds row
func (t *Trie) AddRow(row *Row) {
	t.rows = append(t.rows, row)
}

// AddCmds adds cmd to cmds
func (t *Trie) AddCmds(cmd []rune) {
	t.cmds = append(t.cmds, cmd)
}

// GetLastOnPath returns patch commands
func (t Trie) GetLastOnPath(key []rune) ([]rune, bool) {
	var l []rune // last
	var ok bool

	now, err := t.Row(t.root)
	if err != nil {
		return l, ok
	}

	var w int32
	se := NewStrEnum(key, t.forward)
	for i := 0; i < len(key)-1; i++ {
		c := se.Next()
		w = now.Cmd(c)
		if w >= 0 {
			l = t.cmds[w]
			ok = true
		}

		w = now.Ref(c)
		if w >= 0 {
			now, err = t.Row(w)
			if err != nil {
				return l, ok
			}
		} else {
			return l, ok
		}
	}
	w = now.Cmd(se.Next())
	if w >= 0 {
		return t.cmds[w], true
	}
	return l, ok
}

// Row returns row from trie
func (t Trie) Row(i int32) (*Row, error) {
	if i < 0 || i >= int32(len(t.rows)) {
		err := fmt.Errorf("row %d not found", i)
		return nil, err
	}
	return t.rows[i], nil
}
