// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"unicode"
)

func ignoreSpace(l *lexer) {
	if isSpace(l.next()) {
		l.ignore()
	} else {
		l.backup()
	}
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
