// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleCode(s *StateBlock, startLine, endLine int, _ bool) (_ bool) {
	if s.TShift[startLine]-s.BlkIndent < 4 {
		return
	}

	nextLine := startLine + 1
	last := nextLine

	for nextLine < endLine {
		if s.IsLineEmpty(nextLine) {
			nextLine++
			continue
		}

		if s.TShift[nextLine]-s.BlkIndent > 3 {
			nextLine++
			last = nextLine
			continue
		}

		break
	}

	s.Line = nextLine
	s.PushToken(&CodeBlock{
		Content: s.Lines(startLine, last, 4+s.BlkIndent, true),
		Map:     [2]int{startLine, s.Line},
	})

	return true
}
