// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleStrikeThrough(s *StateInline, silent bool) (_ bool) {
	start := s.Pos
	max := s.PosMax
	src := s.Src

	if src[start] != '~' {
		return
	}

	if silent {
		return
	}

	canOpen, canClose, delims := scanDelims(s, start)
	startCount := delims
	if !canOpen {
		s.Pos += startCount
		s.Pending.WriteString(src[start:s.Pos])
		return true
	}

	stack := startCount / 2
	if stack <= 0 {
		return
	}
	s.Pos = start + startCount

	var found bool
	for s.Pos < max {
		if src[s.Pos] == '~' {
			canOpen, canClose, delims = scanDelims(s, s.Pos)
			count := delims
			tagCount := count / 2
			if canClose {
				if tagCount >= stack {
					s.Pos += count - 2
					found = true
					break
				}
				stack -= tagCount
				s.Pos += count
				continue
			}

			if canOpen {
				stack += tagCount
			}
			s.Pos += count
			continue
		}

		s.Md.Inline.SkipToken(s)
	}

	if !found {
		s.Pos = start
		return
	}

	s.PosMax = s.Pos
	s.Pos = start + 2

	s.PushOpeningToken(&StrikethroughOpen{})

	s.Md.Inline.Tokenize(s)

	s.PushClosingToken(&StrikethroughClose{})

	s.Pos = s.PosMax + 2
	s.PosMax = max

	return true
}
