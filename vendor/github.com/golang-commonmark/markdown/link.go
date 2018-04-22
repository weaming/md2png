// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

func ruleLink(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '[' {
		return
	}

	labelStart := pos + 1
	labelEnd := parseLinkLabel(s, pos, true)
	if labelEnd < 0 {
		return
	}

	var href, title, label string
	oldPos := pos
	pos = labelEnd + 1
	max := s.PosMax
	if pos < max && src[pos] == '(' {
		pos = skipws(src, pos+1, max)
		if pos >= max {
			return
		}

		url, endpos, ok := parseLinkDestination(src, pos, s.PosMax)
		if ok {
			url = normalizeLink(url)
			if validateLink(url) {
				href = url
				pos = endpos
			}
		}

		start := pos
		pos = skipws(src, pos, max)
		if pos >= max {
			return
		}

		title, _, endpos, ok = parseLinkTitle(src, pos, s.PosMax)
		if pos < max && start != pos && ok {
			pos = skipws(src, endpos, max)
		}

		if pos >= max || src[pos] != ')' {
			s.Pos = oldPos
			return
		}

		pos++

	} else {
		if s.Env.References == nil {
			return
		}

		pos = skipws(src, pos, max)

		if pos < max && src[pos] == '[' {
			start := pos + 1
			pos = parseLinkLabel(s, pos, false)
			if pos >= 0 {
				label = src[start:pos]
				pos++
			} else {
				pos = labelEnd + 1
			}
		} else {
			pos = labelEnd + 1
		}

		if label == "" {
			label = src[labelStart:labelEnd]
		}

		ref, ok := s.Env.References[normalizeReference(label)]
		if !ok {
			s.Pos = oldPos
			return
		}

		href = ref["href"]
		title = ref["title"]
	}

	if !silent {
		s.Pos = labelStart
		s.PosMax = labelEnd

		s.PushOpeningToken(&LinkOpen{
			Href:  href,
			Title: title,
		})

		s.Md.Inline.Tokenize(s)

		s.PushClosingToken(&LinkClose{})
	}

	s.Pos = pos
	s.PosMax = max

	return true
}
