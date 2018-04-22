// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

var blockquoteTerminatedBy []BlockRule

func ruleBlockQuote(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	src := s.Src

	if src[pos] != '>' {
		return
	}

	if silent {
		return true
	}

	pos++
	max := s.EMarks[startLine]

	if pos < max && src[pos] == ' ' {
		pos++
	}

	oldIndent := s.BlkIndent
	s.BlkIndent = 0

	oldBMarks := []int{s.BMarks[startLine]}
	s.BMarks[startLine] = pos

	if pos < max {
		pos = s.SkipSpaces(pos)
	}
	lastLineEmpty := pos >= max

	oldTShift := []int{s.TShift[startLine]}
	s.TShift[startLine] = pos - s.BMarks[startLine]

	nextLine := startLine + 1
outer:
	for ; nextLine < endLine; nextLine++ {
		shift := s.TShift[nextLine]
		if shift < oldIndent {
			break
		}

		pos = s.BMarks[nextLine] + shift
		max = s.EMarks[nextLine]

		if pos >= max {
			break
		}

		if src[pos] == '>' {
			pos++
			if pos < max && src[pos] == ' ' {
				pos++
			}

			oldBMarks = append(oldBMarks, s.BMarks[nextLine])
			s.BMarks[nextLine] = pos

			if pos < max {
				pos = s.SkipSpaces(pos)
			}
			lastLineEmpty = pos >= max

			oldTShift = append(oldTShift, s.TShift[nextLine])
			s.TShift[nextLine] = pos - s.BMarks[nextLine]

			continue
		}

		if lastLineEmpty {
			break
		}

		for _, r := range blockquoteTerminatedBy {
			if r(s, nextLine-1, endLine, true) {
				break outer
			}
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}

		oldBMarks = append(oldBMarks, s.BMarks[nextLine])
		oldTShift = append(oldTShift, s.TShift[nextLine])

		s.TShift[nextLine] = -1
	}

	oldParentType := s.ParentType
	s.ParentType = ptBlockQuote
	tok := &BlockquoteOpen{
		Map: [2]int{startLine, 0},
	}
	s.PushOpeningToken(tok)

	s.Md.Block.Tokenize(s, startLine, nextLine)

	s.PushClosingToken(&BlockquoteClose{})
	s.ParentType = oldParentType
	tok.Map[1] = s.Line

	for i := 0; i < len(oldTShift); i++ {
		s.BMarks[startLine+i] = oldBMarks[i]
		s.TShift[startLine+i] = oldTShift[i]
	}
	s.BlkIndent = oldIndent

	return true
}
