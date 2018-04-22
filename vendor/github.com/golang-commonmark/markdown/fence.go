// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

var fence [256]bool

func init() {
	fence['~'], fence['`'] = true, true
}

func ruleFence(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	max := s.EMarks[startLine]
	src := s.Src

	if pos+3 > max {
		return
	}

	marker := src[pos]

	if !fence[marker] {
		return
	}

	mem := pos
	pos = s.SkipBytes(pos, marker)
	len := pos - mem
	if len < 3 {
		return
	}

	params := strings.TrimSpace(src[pos:max])

	if strings.IndexByte(params, '`') >= 0 {
		return
	}

	if silent {
		return true
	}

	nextLine := startLine
	haveEndMarker := false

	for {
		nextLine++
		if nextLine >= endLine {
			break
		}

		mem = s.BMarks[nextLine] + s.TShift[nextLine]
		pos = mem
		max = s.EMarks[nextLine]

		if pos >= max {
			continue
		}

		if s.TShift[nextLine] < s.BlkIndent {
			break
		}

		if src[pos] != marker {
			continue
		}

		if s.TShift[nextLine]-s.BlkIndent > 3 {
			continue
		}

		pos = s.SkipBytes(pos, marker)

		if pos-mem < len {
			continue
		}

		pos = s.SkipSpaces(pos)
		if pos < max {
			continue
		}

		haveEndMarker = true

		break
	}

	s.Line = nextLine
	if haveEndMarker {
		s.Line++
	}

	s.PushToken(&Fence{
		Params:  params,
		Content: s.Lines(startLine+1, nextLine, s.TShift[startLine], true),
		Map:     [2]int{startLine, nextLine},
	})

	return true
}
