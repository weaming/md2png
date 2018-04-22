// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

var referenceTerminatedBy []BlockRule

func ruleReference(s *StateBlock, startLine, _ int, silent bool) (_ bool) {
	pos := s.BMarks[startLine] + s.TShift[startLine]
	src := s.Src

	if src[pos] != '[' {
		return
	}

	pos++
	max := s.EMarks[startLine]

	for pos < max {
		if src[pos] == ']' && src[pos-1] != '\\' {
			if pos+1 == max {
				return
			}
			if src[pos+1] != ':' {
				return
			}
			break
		}
		pos++
	}

	nextLine := startLine + 1
	endLine := s.LineMax
outer:
	for ; nextLine < endLine && !s.IsLineEmpty(nextLine); nextLine++ {
		if s.TShift[nextLine]-s.BlkIndent > 3 {
			continue
		}

		for _, r := range referenceTerminatedBy {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}
	}

	str := strings.TrimSpace(s.Lines(startLine, nextLine, s.BlkIndent, false))
	max = len(str)
	lines := 0
	var labelEnd int
	for pos = 1; pos < max; pos++ {
		b := str[pos]
		if b == '[' {
			return
		} else if b == ']' {
			labelEnd = pos
			break
		} else if b == '\n' {
			lines++
		} else if b == '\\' {
			pos++
			if pos < max && str[pos] == '\n' {
				lines++
			}
		}
	}

	if labelEnd <= 0 || labelEnd+1 >= max || str[labelEnd+1] != ':' {
		return
	}

	for pos = labelEnd + 2; pos < max; pos++ {
		b := str[pos]
		if b == '\n' {
			lines++
		} else if b != ' ' {
			break
		}
	}

	href, endpos, ok := parseLinkDestination(str, pos, max)
	if !ok {
		return
	}
	href = normalizeLink(href)
	if !validateLink(href) {
		return
	}
	pos = endpos

	savedPos := pos
	savedLineNo := lines

	start := pos
	for ; pos < max; pos++ {
		b := str[pos]
		if b == '\n' {
			lines++
		} else if b != ' ' {
			break
		}
	}

	title, nlines, endpos, ok := parseLinkTitle(str, pos, max)
	if pos < max && start != pos && ok {
		pos = endpos
		lines += nlines
	} else {
		pos = savedPos
		lines = savedLineNo
	}

	for pos < max && str[pos] == ' ' {
		pos++
	}

	if pos < max && str[pos] != '\n' {
		if title != "" {
			title = ""
			pos = savedPos
			lines = savedLineNo
			for pos < max && src[pos] == ' ' {
				pos++
			}
		}
	}

	if pos < max && str[pos] != '\n' {
		return
	}

	label := normalizeReference(str[1:labelEnd])
	if label == "" {
		return false
	}

	if silent {
		return true
	}

	if s.Env.References == nil {
		s.Env.References = make(map[string]map[string]string)
	}
	if _, ok := s.Env.References[label]; !ok {
		s.Env.References[label] = map[string]string{
			"title": title,
			"href":  href,
		}
	}

	s.Line = startLine + lines + 1

	return true
}
