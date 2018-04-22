// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import (
	"strconv"

	"github.com/golang-commonmark/markdown/byteutil"
)

var (
	bullet     [256]bool
	afterdigit [256]bool
)

func init() {
	bullet['*'], bullet['+'], bullet['-'] = true, true, true
	afterdigit[')'], afterdigit['.'] = true, true
}

var listTerminatedBy []BlockRule

func skipBulletListMarker(s *StateBlock, startLine int) int {
	pos := s.BMarks[startLine] + s.TShift[startLine]
	src := s.Src

	if !bullet[src[pos]] {
		return -1
	}

	pos++
	max := s.EMarks[startLine]

	if pos < max && src[pos] != ' ' {
		return -1
	}

	return pos
}

func skipOrderedListMarker(s *StateBlock, startLine int) int {
	pos := s.BMarks[startLine] + s.TShift[startLine]
	max := s.EMarks[startLine]

	if pos+1 >= max {
		return -1
	}

	src := s.Src
	b := src[pos]

	if !byteutil.IsDigit(b) {
		return -1
	}

	start := pos

	for {
		if pos >= max {
			return -1
		}

		b = src[pos]
		pos++

		if byteutil.IsDigit(b) {
			if pos-start > 9 {
				return -1
			}
			continue
		}

		if afterdigit[b] {
			break
		}

		return -1
	}

	if pos < max && src[pos] != ' ' {
		return -1
	}

	return pos
}

func markParagraphsTight(s *StateBlock, idx int) {
	level := s.Level + 2
	tokens := s.Tokens

	for i := idx + 2; i < len(tokens)-2; i++ {
		if tokens[i].Level() == level {
			if tok, ok := tokens[i].(*ParagraphOpen); ok {
				tok.Tight = true
				i += 2
				tokens[i].(*ParagraphClose).Tight = true
			}
		}
	}
}

func ruleList(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	shift := s.TShift[startLine]
	if shift < 0 {
		return
	}

	isOrdered := false
	posAfterMarker := skipOrderedListMarker(s, startLine)
	if posAfterMarker > 0 {
		isOrdered = true
	} else {
		posAfterMarker = skipBulletListMarker(s, startLine)
		if posAfterMarker < 0 {
			return
		}
	}

	src := s.Src
	markerChar := src[posAfterMarker-1]

	if silent {
		return true
	}

	tokenIdx := len(s.Tokens)

	var listMap *[2]int
	if isOrdered {
		start := s.BMarks[startLine] + shift
		markerValue, _ := strconv.Atoi(src[start : posAfterMarker-1])

		tok := &OrderedListOpen{
			Order: markerValue,
			Map:   [2]int{startLine, 0},
		}
		s.PushOpeningToken(tok)
		listMap = &tok.Map
	} else {
		tok := &BulletListOpen{
			Map: [2]int{startLine, 0},
		}
		s.PushOpeningToken(tok)
		listMap = &tok.Map
	}

	nextLine := startLine
	prevEmptyEnd := false

	tight := true
outer:
	for nextLine < endLine {
		contentStart := s.SkipSpaces(posAfterMarker)
		max := s.EMarks[nextLine]

		var indentAfterMarker int
		if contentStart >= max {
			indentAfterMarker = 1
		} else {
			indentAfterMarker = contentStart - posAfterMarker
		}

		if indentAfterMarker > 4 {
			indentAfterMarker = 1
		}

		indent := posAfterMarker - s.BMarks[nextLine] + indentAfterMarker

		tok := &ListItemOpen{
			Map: [2]int{startLine, 0},
		}
		s.PushOpeningToken(tok)
		itemMap := &tok.Map

		oldIndent := s.BlkIndent
		oldTight := s.Tight
		oldTShift := s.TShift[startLine]
		oldParentType := s.ParentType
		s.TShift[startLine] = contentStart - s.BMarks[startLine]
		s.BlkIndent = indent
		s.Tight = true
		s.ParentType = ptList

		s.Md.Block.Tokenize(s, startLine, endLine)

		if !s.Tight || prevEmptyEnd {
			tight = false
		}
		prevEmptyEnd = s.Line-startLine > 1 && s.IsLineEmpty(s.Line-1)
		if prevEmptyEnd {
			lastToken := s.Tokens[len(s.Tokens)-1]
			if _, ok := lastToken.(*BlockquoteClose); ok {
				prevEmptyEnd = false
			}
		}

		s.BlkIndent = oldIndent
		s.TShift[startLine] = oldTShift
		s.Tight = oldTight
		s.ParentType = oldParentType

		s.PushClosingToken(&ListItemClose{})

		startLine = s.Line
		nextLine = startLine
		(*itemMap)[1] = nextLine
		contentStart = s.BMarks[startLine]

		if nextLine >= endLine {
			break
		}

		if s.IsLineEmpty(nextLine) {
			break
		}

		if s.TShift[nextLine] < s.BlkIndent {
			break
		}

		for _, r := range listTerminatedBy {
			if r(s, nextLine, endLine, true) {
				break outer
			}
		}

		if isOrdered {
			posAfterMarker = skipOrderedListMarker(s, nextLine)
			if posAfterMarker < 0 {
				break
			}
		} else {
			posAfterMarker = skipBulletListMarker(s, nextLine)
			if posAfterMarker < 0 {
				break
			}
		}

		if markerChar != src[posAfterMarker-1] {
			break
		}
	}

	if isOrdered {
		s.PushClosingToken(&OrderedListClose{})
	} else {
		s.PushClosingToken(&BulletListClose{})
	}
	(*listMap)[1] = nextLine

	s.Line = nextLine

	if tight {
		markParagraphsTight(s, tokenIdx)
	}

	return true
}
