// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

var term [256]bool

func init() {
	for _, b := range "\n!#$%&*+-:<=>@[\\]^_`{}~" {
		term[b] = true
	}
}

func ruleText(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	max := s.PosMax
	src := s.Src

	for pos < max && !term[src[pos]] {
		pos++
	}
	if pos == s.Pos {
		return
	}

	if !silent {
		s.Pending.WriteString(src[s.Pos:pos])
	}

	s.Pos = pos

	return true
}
