// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

type ParserBlock struct {
}

type BlockRule func(*StateBlock, int, int, bool) bool

var blockRules []BlockRule

func (b ParserBlock) Parse(src []byte, md *Markdown, env *Environment) []Token {
	str, bMarks, eMarks, tShift := normalizeAndIndex(src)
	bMarks = append(bMarks, len(str))
	eMarks = append(eMarks, len(str))
	tShift = append(tShift, 0)
	var s StateBlock
	s.BMarks = bMarks
	s.EMarks = eMarks
	s.TShift = tShift
	s.LineMax = len(bMarks) - 1
	s.Src = str
	s.Md = md
	s.Env = env

	b.Tokenize(&s, s.Line, s.LineMax)

	return s.Tokens
}

func (ParserBlock) Tokenize(s *StateBlock, startLine, endLine int) {
	line := startLine
	hasEmptyLines := false
	maxNesting := s.Md.MaxNesting

	for line < endLine {
		line = s.SkipEmptyLines(line)
		s.Line = line
		if line >= endLine {
			break
		}

		if s.TShift[line] < s.BlkIndent {
			break
		}

		if s.Level >= maxNesting {
			s.Line = endLine
			break
		}

		for _, r := range blockRules {
			if r(s, line, endLine, false) {
				break
			}
		}

		s.Tight = !hasEmptyLines

		if s.IsLineEmpty(s.Line - 1) {
			hasEmptyLines = true
		}

		line = s.Line

		if line < endLine && s.IsLineEmpty(line) {
			hasEmptyLines = true
			line++

			if line < endLine && s.ParentType == ptList && s.IsLineEmpty(line) {
				break
			}
			s.Line = line
		}
	}
}
