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

package egothor // import "github.com/kreativka/gostempel/egothor"

import (
	"fmt"

	"github.com/blevesearch/stempel/javadata"
)

// Trie implements egothor trie
type Trie struct {
	cmds    [][]rune
	forward bool
	root    int32
	rows    []*row
}

// NewTrie returns egothor trie
func NewTrie(r *javadata.Reader) (*Trie, error) {
	var err error
	rv := &Trie{}

	rv.forward, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("error reading forward value: %v", err)
	}

	rv.root, err = r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("error reading root value: %v", err)
	}

	i, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("error reading number of cmds")
	}
	for i > 0 {
		cmd, err := r.ReadUTF()
		if err != nil {
			return nil, fmt.Errorf("error reading command: %v", err)
		}
		rv.cmds = append(rv.cmds, []rune(cmd))
		i--
	}

	i, err = r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("error reading number of rows: %v", err)
	}
	for i > 0 {
		row, err := newRow(r)
		if err != nil {
			return nil, fmt.Errorf("error reading row: %v", err)
		}

		rv.rows = append(rv.rows, row)
		i--
	}
	return rv, nil
}

// GetLastOnPath returns patch commands to apply by DiffApply to get stemmed
// token.
// This walks rune by rune in key through rows, checking for commands. When
// it reaches end of string, return commands from this row or if empty return
// last set of commands.
func (t *Trie) GetLastOnPath(key []rune) []rune {
	var last []rune
	var w int32
	now := t.row(t.root)
	e := newStrEnum(key, t.forward)
	for i := 0; i < len(key)-1; i++ {
		r, err := e.next()
		if err != nil {
			return last
		}
		w = now.cmd(r)
		if w >= 0 {
			last = t.cmds[w]
		}

		w = now.ref(r)
		if w >= 0 {
			now = t.row(w)
		} else {
			return last
		}
	}
	r, err := e.next()
	if err != nil {
		return last
	}
	w = now.cmd(r)
	if w >= 0 {
		return t.cmds[w]
	}
	return last
}

func (t *Trie) row(i int32) *row {
	if i < 0 || i >= int32(len(t.rows)) {
		return nil
	}
	return t.rows[i]
}
