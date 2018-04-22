// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "strings"

var hdr [256]bool

func init() {
	for _, b := range "-:| " {
		hdr[b] = true
	}
}

func getLine(s *StateBlock, line int) string {
	pos := s.BMarks[line] + s.BlkIndent
	max := s.EMarks[line]
	if pos >= max {
		return ""
	}
	return s.Src[pos:max]
}

func escapedSplit(s string) (result []string) {
	escapes := 0
	lastPos := 0
	backticked := false
	lastBackTick := 0
	pos := 0

	if len(s) > 0 && s[len(s)-1] == '|' {
		s = s[:len(s)-1]
	}
	if len(s) > 0 && s[0] == '|' {
		pos++
		lastPos++
	}

	for pos < len(s) {
		b := s[pos]
		switch {
		case b == '`' && escapes%2 == 0:
			backticked = !backticked
			lastBackTick = pos
		case b == '|' && escapes%2 == 0 && !backticked:
			result = append(result, s[lastPos:pos])
			lastPos = pos + 1
		case b == '\\':
			escapes++
		default:
			escapes = 0
		}

		pos++

		if pos == len(s) && backticked {
			backticked = false
			pos = lastBackTick + 1
		}
	}

	result = append(result, s[lastPos:])

	return
}

func ruleTable(s *StateBlock, startLine, endLine int, silent bool) (_ bool) {
	if !s.Md.Tables {
		return
	}

	if startLine+2 > endLine {
		return
	}

	nextLine := startLine + 1

	if s.TShift[nextLine] < s.BlkIndent {
		return
	}

	pos := s.BMarks[nextLine] + s.TShift[nextLine]
	if pos >= s.EMarks[nextLine] {
		return
	}

	src := s.Src
	if !hdr[src[pos]] {
		return
	}

	lineText := getLine(s, startLine+1)
	if !isHeaderLine(lineText) {
		return
	}

	rows := strings.Split(lineText, "|")
	if len(rows) < 2 {
		return
	}
	var aligns []Align
	for i := 0; i < len(rows); i++ {
		t := strings.TrimSpace(rows[i])
		if t == "" {
			continue
		}

		if t[len(t)-1] == ':' {
			if t[0] == ':' {
				aligns = append(aligns, AlignCenter)
			} else {
				aligns = append(aligns, AlignRight)
			}
		} else if t[0] == ':' {
			aligns = append(aligns, AlignLeft)
		} else {
			aligns = append(aligns, AlignNone)
		}
	}

	lineText = strings.TrimSpace(getLine(s, startLine))
	if strings.IndexByte(lineText, '|') == -1 {
		return
	}

	rows = escapedSplit(lineText)
	if len(aligns) != len(rows) {
		return
	}

	if silent {
		return true
	}

	tableTok := &TableOpen{
		Map: [2]int{startLine, 0},
	}
	s.PushOpeningToken(tableTok)
	s.PushOpeningToken(&TheadOpen{
		Map: [2]int{startLine, startLine + 1},
	})
	s.PushOpeningToken(&TrOpen{
		Map: [2]int{startLine, startLine + 1},
	})

	for i := 0; i < len(rows); i++ {
		s.PushOpeningToken(&ThOpen{
			Align: aligns[i],
			Map:   [2]int{startLine, startLine + 1},
		})
		s.PushToken(&Inline{
			Content: strings.TrimSpace(rows[i]),
			Map:     [2]int{startLine, startLine + 1},
		})
		s.PushClosingToken(&ThClose{})
	}

	s.PushClosingToken(&TrClose{})
	s.PushClosingToken(&TheadClose{})

	tbodyTok := &TbodyOpen{
		Map: [2]int{startLine + 2, 0},
	}
	s.PushOpeningToken(tbodyTok)

	for nextLine = startLine + 2; nextLine < endLine; nextLine++ {
		shift := s.TShift[nextLine]
		if shift >= 0 && shift < s.BlkIndent {
			break
		}

		lineText = strings.TrimSpace(getLine(s, nextLine))
		if strings.IndexByte(lineText, '|') == -1 {
			break
		}
		rows = escapedSplit(lineText)
		if len(rows) < len(aligns) {
			rows = append(rows, make([]string, len(aligns)-len(rows))...)
		} else if len(rows) > len(aligns) {
			rows = rows[:len(aligns)]
		}

		s.PushOpeningToken(&TrOpen{})
		for i := 0; i < len(rows); i++ {
			s.PushOpeningToken(&TdOpen{Align: aligns[i]})
			s.PushToken(&Inline{
				Content: strings.TrimSpace(rows[i]),
			})
			s.PushClosingToken(&TdClose{})
		}
		s.PushClosingToken(&TrClose{})
	}

	s.PushClosingToken(&TbodyClose{})
	s.PushClosingToken(&TableClose{})

	tableTok.Map[1] = nextLine
	tbodyTok.Map[1] = nextLine
	s.Line = nextLine

	return true
}
