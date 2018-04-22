// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import (
	"strings"
	"unicode"
)

type stackItem struct {
	token  int
	text   []rune
	pos    int
	single bool
	level  int
}

func nextQuoteIndex(s []rune, from int) int {
	for i := from; i < len(s); i++ {
		r := s[i]
		if r == '\'' || r == '"' {
			return i
		}
	}
	return -1
}

func replaceQuotes(tokens []Token, s *StateCore) {
	var stack []stackItem
	var changed map[int][]rune

	for i, tok := range tokens {
		thisLevel := tok.Level()

		j := len(stack) - 1
		for j >= 0 {
			if stack[j].level <= thisLevel {
				break
			}
			j--
		}
		stack = stack[:j+1]

		if tok, ok := tok.(*Text); ok {
			if !strings.ContainsAny(tok.Content, `"'`) {
				continue
			}

			text := []rune(tok.Content)
			pos := 0
			max := len(text)

		loop:
			for pos < max {
				index := nextQuoteIndex(text, pos)
				if index < 0 {
					break
				}

				canOpen := true
				canClose := true
				pos = index + 1
				isSingle := text[index] == '\''

				var lastChar rune
				if index > 0 {
					lastChar = text[index-1]
				}
				var nextChar rune
				if pos < max {
					nextChar = text[pos]
				}

				isLastSpaceOrStart := index == 0 || unicode.IsSpace(lastChar)
				isNextSpaceOrEnd := pos == max || unicode.IsSpace(nextChar)
				isLastPunct := !isLastSpaceOrStart && (isMarkdownPunct(lastChar) || unicode.IsPunct(lastChar))
				isNextPunct := !isNextSpaceOrEnd && (isMarkdownPunct(nextChar) || unicode.IsPunct(nextChar))

				if isNextSpaceOrEnd {
					canOpen = false
				} else if isNextPunct {
					if !(isLastSpaceOrStart || isLastPunct) {
						canOpen = false
					}
				}

				if isLastSpaceOrStart {
					canClose = false
				} else if isLastPunct {
					if !(isNextSpaceOrEnd || isNextPunct) {
						canClose = false
					}
				}

				if nextChar == '"' && text[index] == '"' {
					if lastChar >= '0' && lastChar <= '9' {
						canClose = false
						canOpen = false
					}
				}

				if canOpen && canClose {
					canOpen = false
					canClose = isNextPunct
				}

				if !canOpen && !canClose {
					if isSingle {
						text[index] = '’'
						if changed == nil {
							changed = make(map[int][]rune)
						}
						if _, ok := changed[i]; !ok {
							changed[i] = text
						}
					}
					continue
				}

				if canClose {
					for j := len(stack) - 1; j >= 0; j-- {
						item := stack[j]
						if item.level < thisLevel {
							break
						}
						if item.single == isSingle && item.level == thisLevel {
							if changed == nil {
								changed = make(map[int][]rune)
							}
							if isSingle {
								item.text[item.pos] = s.Md.options.Quotes[2]
								text[index] = s.Md.options.Quotes[3]
							} else {
								item.text[item.pos] = s.Md.options.Quotes[0]
								text[index] = s.Md.options.Quotes[1]
							}
							if _, ok := changed[i]; !ok {
								changed[i] = text
							}
							if ii := item.token; ii != i {
								if _, ok := changed[ii]; !ok {
									changed[ii] = item.text
								}
							}
							stack = stack[:j]
							continue loop
						}
					}
				}

				if canOpen {
					stack = append(stack, stackItem{
						token:  i,
						text:   text,
						pos:    index,
						single: isSingle,
						level:  thisLevel,
					})
				} else if canClose && isSingle {
					text[index] = '’'
					if changed == nil {
						changed = make(map[int][]rune)
					}
					if _, ok := changed[i]; !ok {
						changed[i] = text
					}
				}
			}
		}
	}

	if changed != nil {
		for i, text := range changed {
			tokens[i].(*Text).Content = string(text)
		}
	}
}

func ruleSmartQuotes(s *StateCore) {
	if !s.Md.Typographer {
		return
	}

	tokens := s.Tokens
	for i := len(tokens) - 1; i >= 0; i-- {
		tok := tokens[i]
		if tok, ok := tok.(*Inline); ok {
			replaceQuotes(tok.Children, s)
		}
	}
}
