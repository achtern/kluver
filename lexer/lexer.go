// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

type stateFn func(*lexer) stateFn

func Lex(name, input string) (*lexer, chan token) {
	l := &lexer{
		name:   name,
		input:  input,
		tokens: make(chan token),
	}

	go l.run()
	return l, l.tokens
}

func (l *lexer) run() {
	for state := lexVoid; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) hasPrefix(pre string) bool {
	return strings.HasPrefix(l.input[l.pos:], pre)
}

func (l *lexer) testPrefix(pre string, token tokenType) bool {
	if l.hasPrefix(pre) {
		if l.pos > l.start {
			l.emit(token)
		}

		return true
	}

	return false
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) isNumber() bool {
	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}

	l.accept("i")

	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

func (l *lexer) lexStatement(keyword string, emit tokenType, next stateFn) stateFn {
	l.pos += len(keyword)
	l.emit(emit)
	return next
}
