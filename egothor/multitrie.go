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

	"github.com/blevesearch/stempel/javadata"
)

// End of message
const eom = '*'

// MultiTrie implements egothor MultiTrie
type MultiTrie struct {
	by      int32
	forward bool
	tries   []*Trie
}

// NewMultiTrie returns egothor MultiTrie
func NewMultiTrie(r *javadata.Reader) (*MultiTrie, error) {
	var err error
	rv := &MultiTrie{}

	rv.forward, err = r.ReadBool()
	if err != nil {
		return nil, fmt.Errorf("error reading forward value: %v", err)
	}

	rv.by, err = r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("error reading by value: %v", err)
	}

	i, err := r.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("error reading number of tries: %v", err)
	}
	for i > 0 {
		trie, err := NewTrie(r)
		if err != nil {
			return nil, fmt.Errorf("error reading new trie: %v", err)
		}

		rv.tries = append(rv.tries, trie)
		i--
	}
	return rv, nil
}

// GetLastOnPath returns patch commands.
func (t *MultiTrie) GetLastOnPath(key []rune) []rune {
	lastKey := key
	lastR := ' '
	p := make(map[int][]rune)
	var rv []rune

	for i := 0; i < len(t.tries); i++ {
		r := t.tries[i].GetLastOnPath(lastKey)
		if len(r) == 0 || len(r) == 1 && r[0] == eom {
			return rv
		}
		if cannotFollow(lastR, r[0]) {
			return rv
		}
		lastR = r[len(r)-2]
		p[i] = r
		if p[i][0] == '-' {
			if i > 0 {
				var err error
				key, err = t.skip(key, lengthPP(p[i-1]))
				if err != nil {
					return rv
				}
			}
			var err error
			key, err = t.skip(key, lengthPP(p[i]))
			if err != nil {
				return rv
			}
		}
		rv = append(rv, r...)
		if len(key) != 0 {
			lastKey = key
		}
	}
	return rv
}

func (t *MultiTrie) skip(in []rune, count int) ([]rune, error) {
	if count > len(in) {
		return nil, fmt.Errorf("index out of bounds")
	}
	if t.forward {
		return in[count:], nil
	}
	return in[:len(in)-count], nil
}

func cannotFollow(after, goes rune) bool {
	if after == '-' || after == 'D' {
		return after == goes
	}
	return false
}

func lengthPP(cmd []rune) int {
	rv := 0
	for i := 0; i < len(cmd); i++ {
		switch cmd[i] {
		case '-', 'D':
			i++
			rv += int(cmd[i] - 'a' + 1)
		case 'R':
			i++
			rv++
			fallthrough
		case 'I':
		}
	}
	return rv
}
