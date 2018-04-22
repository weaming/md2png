// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

func ruleAutolink(s *StateInline, silent bool) (_ bool) {
	pos := s.Pos
	src := s.Src

	if src[pos] != '<' {
		return
	}

	tail := src[pos:]

	if strings.IndexByte(tail, '>') < 0 {
		return
	}

	link := matchAutolink(tail)
	if link != "" {
		href := normalizeLink(link)
		if !validateLink(href) {
			return
		}

		if !silent {
			s.PushOpeningToken(&LinkOpen{Href: href})
			s.PushToken(&Text{Content: normalizeLinkText(link)})
			s.PushClosingToken(&LinkClose{})
		}

		s.Pos += len(link) + 2

		return true
	}

	email := matchEmail(tail)
	if email != "" {
		href := normalizeLink("mailto:" + email)
		if !validateLink(href) {
			return
		}

		if !silent {
			s.PushOpeningToken(&LinkOpen{Href: href})
			s.PushToken(&Text{Content: email})
			s.PushClosingToken(&LinkClose{})
		}

		s.Pos += len(email) + 2

		return true
	}

	return
}
