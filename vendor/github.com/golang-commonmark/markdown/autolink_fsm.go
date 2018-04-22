// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "github.com/golang-commonmark/markdown/byteutil"

var (
	schemecs [256]bool
	linkcs   [256]bool
	emailcs  [256]bool
)

func init() {
	for _, b := range "-.abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		schemecs[b] = true
	}
	for i := 0x21; i <= 0xff; i++ {
		if !(i == '<' || i == '>') {
			linkcs[i] = true
		}
	}
	for _, b := range "+-._0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		emailcs[b] = true
	}
}

func matchAutolink(s string) string {
	if len(s) < 6 || s[0] != '<' {
		return ""
	}

	st := 0
	n := 0
	for i := 1; i < len(s); i++ {
		b := s[i]
		switch st {
		case 0: // initial state
			switch {
			case schemecs[b]:
				st = 1
				n++
			default:
				return ""
			}

		case 1: // h
			switch {
			case schemecs[b]:
				n++
				if n > 23 {
					return ""
				}
			case b == ':':
				schema := byteutil.ToLower(s[1:i])
				if !matchSchema(schema) {
					return ""
				}
				st = 2
			default:
				return ""
			}

		case 2: // http:
			switch {
			case linkcs[b]:
				st = 3
				break
			default:
				return ""
			}
		case 3: // http:/
			switch {
			case linkcs[b]:
				break
			case b == '>':
				return s[1:i]
			default:
				return ""
			}
		}
	}
	return ""
}

func matchEmail(s string) string {
	if len(s) < 8 || s[0] != '<' {
		return ""
	}

	st := 0
	for i := 1; i < len(s); i++ {
		b := s[i]
		switch st {
		case 0: // initial state
			switch {
			case emailcs[b]:
				st = 1
			default:
				return ""
			}

		case 1: // r
			switch {
			case emailcs[b]:
				break
			case b == '@':
				st = 2
			default:
				return ""
			}

		case 2: // root@
			switch {
			case emailcs[b]:
				st = 3
			default:
				return ""
			}

		case 3: // root@l
			switch {
			case emailcs[b]:
				break
			case b == '>':
				return s[1:i]
			default:
				return ""
			}
		}
	}
	return ""
}
