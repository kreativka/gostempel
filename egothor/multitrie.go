package egothor

/*
                    Egothor Software License version 1.00
                    Copyright (C) 1997-2004 Leo Galambos.
                 Copyright (C) 2002-2004 "Egothor developers"
                      on behalf of the Egothor Project.
                             All rights reserved.

   This  software  is  copyrighted  by  the "Egothor developers". If this
   license applies to a single file or document, the "Egothor developers"
   are the people or entities mentioned as copyright holders in that file
   or  document.  If  this  license  applies  to the Egothor project as a
   whole,  the  copyright holders are the people or entities mentioned in
   the  file CREDITS. This file can be found in the same location as this
   license in the distribution.

   Redistribution  and  use  in  source and binary forms, with or without
   modification, are permitted provided that the following conditions are
   met:
    1. Redistributions  of  source  code  must retain the above copyright
       notice, the list of contributors, this list of conditions, and the
       following disclaimer.
    2. Redistributions  in binary form must reproduce the above copyright
       notice, the list of contributors, this list of conditions, and the
       disclaimer  that  follows  these  conditions  in the documentation
       and/or other materials provided with the distribution.
    3. The name "Egothor" must not be used to endorse or promote products
       derived  from  this software without prior written permission. For
       written permission, please contact Leo.G@seznam.cz
    4. Products  derived  from this software may not be called "Egothor",
       nor  may  "Egothor"  appear  in  their name, without prior written
       permission from Leo.G@seznam.cz.

   In addition, we request that you include in the end-user documentation
   provided  with  the  redistribution  and/or  in the software itself an
   acknowledgement equivalent to the following:
   "This product includes software developed by the Egothor Project.
    http://egothor.sf.net/"

   THIS  SOFTWARE  IS  PROVIDED  ``AS  IS''  AND ANY EXPRESSED OR IMPLIED
   WARRANTIES,  INCLUDING,  BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
   MERCHANTABILITY  AND  FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
   IN  NO  EVENT  SHALL THE EGOTHOR PROJECT OR ITS CONTRIBUTORS BE LIABLE
   FOR   ANY   DIRECT,   INDIRECT,  INCIDENTAL,  SPECIAL,  EXEMPLARY,  OR
   CONSEQUENTIAL  DAMAGES  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
   SUBSTITUTE  GOODS  OR  SERVICES;  LOSS  OF  USE,  DATA, OR PROFITS; OR
   BUSINESS  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
   WHETHER  IN  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
   OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
   IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

   This  software  consists  of  voluntary  contributions  made  by  many
   individuals  on  behalf  of  the  Egothor  Project  and was originally
   created by Leo Galambos (Leo.G@seznam.cz).
*/

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
