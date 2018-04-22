// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

var hr [256]bool

func init() {
	hr['*'], hr['-'], hr['_'] = true, true, true
}

func ruleHR(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	pos := s.BMarks[startLine] + shift
	src := s.Src

	marker := src[pos]

	if !hr[marker] {
		return
	}

	pos++
	max := s.EMarks[startLine]

	count := 1
	for pos < max {
		c := src[pos]
		pos++
		if c != marker && c != ' ' {
			return
		}
		if c == marker {
			count++
		}
	}

	if count < 3 {
		return
	}

	if silent {
		return true
	}

	s.Line = startLine + 1
	s.PushToken(&Hr{
		Map: [2]int{startLine, s.Line},
	})

	return true
}
