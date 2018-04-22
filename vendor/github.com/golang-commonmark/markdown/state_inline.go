// Copyright 2015 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package markdown

import "bytes"

type StateInline struct {
	StateCore

	Pos          int
	PosMax       int
	Level        int
	Pending      bytes.Buffer
	PendingLevel int

	Cache map[int]int
}

func (s *StateInline) PushToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	tok.SetLevel(s.Level)
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushOpeningToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	tok.SetLevel(s.Level)
	s.Level++
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushClosingToken(tok Token) {
	if s.Pending.Len() > 0 {
		s.PushPending()
	}
	s.Level--
	tok.SetLevel(s.Level)
	s.PendingLevel = s.Level
	s.Tokens = append(s.Tokens, tok)
}

func (s *StateInline) PushPending() {
	s.Tokens = append(s.Tokens, &Text{
		Content: s.Pending.String(),
		Lvl:     s.PendingLevel,
	})
	s.Pending.Reset()
}
