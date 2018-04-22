// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleNewline(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '\n' {
		return
	}

	pending := s.Pending.Bytes()
	n := len(pending) - 1

	if !silent {
		if n >= 0 && pending[n] == ' ' {
			if n >= 1 && pending[n-1] == ' ' {
				n -= 2
				for n >= 0 && pending[n] == ' ' {
					n--
				}
				s.Pending.Truncate(n + 1)
				s.PushToken(&Hardbreak{})
			} else {
				s.Pending.Truncate(n)
				s.PushToken(&Softbreak{})
			}
		} else {
			s.PushToken(&Softbreak{})
		}
	}

	pos++
	max := s.PosMax

	for pos < max && src[pos] == ' ' {
		pos++
	}

	s.Pos = pos

	return true
}
