// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"github.com/achtern/kluver/lexer"
)

func (slice Tokens) findFirst(value lexer.TokenType) int {
	for p, v := range slice {
		if v.Typ == value {
			return p
		}
	}
	return -1
}

func (slice Tokens) findLast(value lexer.TokenType) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i].Typ == value {
			return i
		}
	}
	return -1
}