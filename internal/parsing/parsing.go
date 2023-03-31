package parsing

import (
	"errors"
	"strings"
	"unicode/utf8"
)

var errEndOfInput = errors.New("end of input reached before lexing finished")

type state struct {
	input  string
	offset int
}

func newState(input string) state {
	return state{
		input:  input,
		offset: 0,
	}
}

func (s state) nextRune() (rune, state) {
	r, w := utf8.DecodeRuneInString(s.remaining())
	return r, s.advance(w)
}

func (s state) remaining() string {
	return s.input[s.offset:]
}

func (s state) advance(n int) state {
	s.offset += n
	return s
}

func (s state) endOfInput() bool {
	return s.offset >= len(s.input)
}

func consumeIf(s state, condition func(rune) bool) (state, bool) {
	r, next := s.nextRune()
	if condition(r) {
		return next, true
	}
	return s, false
}

func consumeWhile(s state, condition func(rune) bool) state {
	for !s.endOfInput() {
		r, next := s.nextRune()
		if !condition(r) {
			break
		}
		s = next
	}
	return s
}

func lexUntil(s state, until rune) (string, state, error) {
	var b strings.Builder
	for !s.endOfInput() {
		var r rune
		r, s = s.nextRune()
		if r == until {
			return b.String(), s, nil
		}
		b.WriteRune(r)
	}
	return b.String(), s, errEndOfInput
}

func lexBalanced(s state, opening, closing rune) (string, state, error) {
	var b strings.Builder
	balance := 0
	for !s.endOfInput() {
		var r rune
		r, s = s.nextRune()
		switch r {
		case closing:
			if balance == 0 {
				return b.String(), s, nil
			}
			balance--
		case opening:
			balance++
		}
		b.WriteRune(r)
	}
	return b.String(), s, errEndOfInput
}
