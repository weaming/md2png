// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

var escaped = make([]bool, 256)

func init() {
	for _, b := range "\\!\"#$%&'()*+,./:;<=>?@[]^_`{|}~-" {
		escaped[b] = true
	}
}

func ruleEscape(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '\\' {
		return
	}

	pos++
	max := s.PosMax

	if pos < max {
		b := src[pos]

		if b < 0x7f && escaped[b] {
			if !silent {
				s.Pending.WriteByte(b)
			}
			s.Pos += 2
			return true
		}

		if b == '\n' {
			if !silent {
				s.PushToken(&Hardbreak{})
			}

			pos++

			for pos < max && src[pos] == ' ' {
				pos++
			}

			s.Pos = pos
			return true
		}
	}

	if !silent {
		s.Pending.WriteByte('\\')
	}

	s.Pos++

	return true
}
