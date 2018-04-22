// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "unicode/utf8"

type ParserInline struct {
}

type InlineRule func(*StateInline, bool) bool

var inlineRules []InlineRule

func (i ParserInline) Parse(src string, md *Markdown, env *Environment) []Token {
	if src == "" {
		return nil
	}

	var s StateInline
	s.Src = src
	s.Md = md
	s.Env = env
	s.PosMax = len(src)
	s.Tokens = s.TokArr[:0]

	i.Tokenize(&s)

	return s.Tokens
}

func (ParserInline) Tokenize(s *StateInline) {
	max := s.PosMax
	src := s.Src
	maxNesting := s.Md.MaxNesting

outer:
	for s.Pos < max {
		if s.Level < maxNesting {
			for _, rule := range inlineRules {
				if rule(s, false) {
					if s.Pos >= max {
						break outer
					}
					continue outer
				}
			}
		}

		r, size := utf8.DecodeRuneInString(src[s.Pos:])
		s.Pending.WriteRune(r)
		s.Pos += size
	}

	if s.Pending.Len() > 0 {
		s.PushPending()
	}
}

func (ParserInline) SkipToken(s *StateInline) {
	pos := s.Pos
	if s.Cache != nil {
		if pos, ok := s.Cache[pos]; ok {
			s.Pos = pos
			return
		}
	} else {
		s.Cache = make(map[int]int)
	}

	if s.Level < s.Md.MaxNesting {
		for _, r := range inlineRules {
			if r(s, true) {
				s.Cache[pos] = s.Pos
				return
			}
		}
	}

	s.Pos++
	s.Cache[pos] = s.Pos
}
