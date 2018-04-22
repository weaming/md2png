// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleBackticks(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '`' {
		return
	}

	start := pos
	pos++
	max := s.PosMax

	for pos < max && src[pos] == '`' {
		pos++
	}

	marker := src[start:pos]

	end := pos

	for {
		for start = end; start < max && src[start] != '`'; start++ {
			// do nothing
		}
		if start >= max {
			break
		}
		end = start + 1

		for end < max && src[end] == '`' {
			end++
		}

		if end-start == len(marker) {
			if !silent {
				s.PushToken(&CodeInline{
					Content: normalizeInlineCode(src[pos:start]),
				})
			}
			s.Pos = end
			return true
		}
	}

	if !silent {
		s.Pending.WriteString(marker)
	}

	s.Pos += len(marker)

	return true
}
