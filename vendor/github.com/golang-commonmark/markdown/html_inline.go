// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleHTMLInline(s *StateInline, silent bool) (_ bool) {
	if !s.Md.HTML {
		return
	}

	pos := s.Pos
	src := s.Src
	if pos+2 >= s.PosMax || src[pos] != '<' {
		return
	}

	if !htmlSecond[src[pos+1]] {
		return
	}

	match := matchHTML(src[pos:])
	if match == "" {
		return
	}

	if !silent {
		s.PushToken(&HTMLInline{Content: match})
	}

	s.Pos += len(match)

	return true
}
