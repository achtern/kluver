// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package builder

import (
	"bytes"
	"github.com/achtern/kluver/lexer"
)

func Build(tokenStream <-chan lexer.Token) string {
	tokens := make([]lexer.Token, 0)

	var buffer bytes.Buffer

	for token := range tokenStream {
		if token.Typ == lexer.TokenVoid {
			continue
		}
		tokens = append(tokens, token)
	}

	for _, token := range tokens {
		buffer.WriteString(token.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}
