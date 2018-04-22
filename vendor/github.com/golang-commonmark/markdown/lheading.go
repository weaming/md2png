// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

var under [256]bool

func init() {
	under['-'], under['='] = true, true
}

func ruleLHeading(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	nextLine := startLine + 1

	if nextLine >= endLine {
		return
	}

	shift := s.TShift[nextLine]
	if shift < s.BlkIndent {
		return
	}

	if shift-s.BlkIndent > 3 {
		return
	}

	pos := s.BMarks[nextLine] + shift
	max := s.EMarks[nextLine]

	if pos >= max {
		return
	}

	src := s.Src
	marker := src[pos]

	if !under[marker] {
		return
	}

	pos = s.SkipBytes(pos, marker)

	pos = s.SkipSpaces(pos)

	if pos < max {
		return
	}

	pos = s.BMarks[startLine] + s.TShift[startLine]

	s.Line = nextLine + 1

	hLevel := 1
	if marker == '-' {
		hLevel++
	}

	s.PushOpeningToken(&HeadingOpen{
		HLevel: hLevel,
		Map:    [2]int{startLine, s.Line},
	})
	s.PushToken(&Inline{
		Content: strings.TrimSpace(src[pos:s.EMarks[startLine]]),
		Map:     [2]int{startLine, s.Line - 1},
	})
	s.PushClosingToken(&HeadingClose{HLevel: hLevel})

	return true
}
