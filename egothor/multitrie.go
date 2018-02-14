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
	"io"

	"github.com/kreativka/gostempel/javaread"
)

// MultiTrie struct
type MultiTrie struct {
	by      int32
	eom     rune
	forward bool
	root    int32
	tries   []*Trie
}

// NewMultiTrie returns MultiTrie
func NewMultiTrie(in io.Reader) *MultiTrie {
	var mt MultiTrie

	mt.SetEom()
	mt.SetForward(javaread.Bool(in))
	mt.SetBY(javaread.Int(in))

	for i := javaread.Int(in); i > 0; i-- {
		mt.Add(NewTrie(in))
	}
	return &mt
}

// Tries return tries
func (m *MultiTrie) Tries() []*Trie {
	return m.tries
}

// Add appends trie to tries
func (m *MultiTrie) Add(trie *Trie) {
	m.tries = append(m.tries, trie)
}

// SetBY sets BY
func (m *MultiTrie) SetBY(by int32) {
	m.by = by
}

// SetEom sets EOM
func (m *MultiTrie) SetEom() {
	m.eom = '*'
}

// SetForward sets forward
func (m *MultiTrie) SetForward(forward bool) {
	m.forward = forward
}

// SetRoot sets root index of trie
func (m *MultiTrie) SetRoot(root int32) {
	m.root = root
}

// GetLastOnPath returns patch commands
func (m *MultiTrie) GetLastOnPath(key string) ([]rune, bool) {
	var res []rune            // Result
	lk := key                 // Last Key
	lc := rune(' ')           // Last char
	p := make(map[int][]rune) // Patch commands

	for i := 0; i < len(m.Tries()); i++ {
		r, ok := m.tries[i].GetLastOnPath(lk)
		if !ok || (r[0] == m.eom && len(r) == 1) {
			return res, true
		}

		if cannotFollow(lc, r[0]) {
			return res, false
		}
		lc = r[len(r)-2]

		p[i] = r
		if p[i][0] == '-' {
			if i > 0 {
				if key, ok = m.skip(key, lenPP(p[i-1])); !ok {
					return res, true
				}
			}
			if key, ok = m.skip(key, lenPP(p[i])); !ok {
				return res, true
			}
		}
		res = append(res, r...)

		if len(key) != 0 {
			lk = key
		}
	}
	return res, true
}

func cannotFollow(after, goes rune) bool {
	if after == '-' || after == 'D' {
		return after == goes
	}
	return false
}

func (m MultiTrie) skip(in string, count int) (string, bool) {
	runes := []rune(in)

	if len(runes)-count < 0 {
		return "", false
	}

	if m.forward {
		return string(runes[count:]), true
	}
	return string(runes[:len(runes)-count]), true
}

func lenPP(cmd []rune) int {
	l := 0 // length
	for i, c := range cmd {
		i++
		switch c {
		case '-':
			l += int(cmd[i]) - 'a' + 1
		case 'D':
			l += int(cmd[i]) - 'a' + 1
		case 'R':
			l++
		}
	}
	return l
}
