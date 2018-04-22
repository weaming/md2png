// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

const (
	ptRoot = iota
	ptList
	ptBlockQuote
)

type StateBlock struct {
	StateCore

	BMarks     []int // offsets of the line beginnings
	EMarks     []int // offsets of the line endings
	TShift     []int // indents for each line
	BlkIndent  int   // required block content indent (in a list etc.)
	Line       int   // line index in the source string
	LineMax    int   // number of lines
	Tight      bool  // loose or tight mode for lists
	ParentType byte  // parent block type
	Level      int
}

func (s *StateBlock) IsLineEmpty(n int) bool {
	return s.BMarks[n]+s.TShift[n] >= s.EMarks[n]
}

func (s *StateBlock) SkipEmptyLines(from int) int {
	for from < s.LineMax && s.IsLineEmpty(from) {
		from++
	}
	return from
}

func (s *StateBlock) SkipSpaces(pos int) int {
	src := s.Src
	for pos < len(src) && src[pos] == ' ' {
		pos++
	}
	return pos
}

func (s *StateBlock) SkipBytes(pos int, b byte) int {
	src := s.Src
	for pos < len(src) && src[pos] == b {
		pos++
	}
	return pos
}

func (s *StateBlock) SkipBytesBack(pos int, b byte, min int) int {
	for pos > min {
		pos--
		if s.Src[pos] != b {
			return pos + 1
		}
	}
	return pos
}

func (s *StateBlock) Lines(begin, end, indent int, keepLastLf bool) string {
	if begin == end {
		return ""
	}

	src := s.Src

	if begin+1 == end {
		shift := s.TShift[begin]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.BMarks[begin] + shift

		last := s.EMarks[begin]
		if keepLastLf && last < len(src) {
			last++
		}

		return src[first:last]
	}

	size := 0
	var firstFirst int
	var previousLast int
	adjoin := true
	for line := begin; line < end; line++ {
		shift := s.TShift[line]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.BMarks[line] + shift
		last := s.EMarks[line]
		if line+1 < end || (keepLastLf && last < len(src)) {
			last++
		}
		size += last - first
		if line == begin {
			firstFirst = first
		} else if previousLast != first {
			adjoin = false
		}
		previousLast = last
	}

	if adjoin {
		return src[firstFirst:previousLast]
	}

	buf := make([]byte, size)
	i := 0
	for line := begin; line < end; line++ {
		shift := s.TShift[line]
		if shift < 0 {
			shift = 0
		} else if shift > indent {
			shift = indent
		}
		first := s.BMarks[line] + shift
		last := s.EMarks[line]
		if line+1 < end || (keepLastLf && last < len(src)) {
			last++
		}

		i += copy(buf[i:], src[first:last])
	}

	return string(buf)
}

func (s *StateBlock) PushToken(tok Token) {
	tok.SetLevel(s.Level)
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateBlock) PushOpeningToken(tok Token) {
	tok.SetLevel(s.Level)
	s.Level++
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateBlock) PushClosingToken(tok Token) {
	s.Level--
	tok.SetLevel(s.Level)
	s.Tokens = append(s.Tokens, tok)
}
