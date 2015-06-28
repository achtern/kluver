// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"bytes"
	"fmt"
	"github.com/achtern/kluver/lexer"
	"strconv"
)

type Shader struct {
	version  int
	vertex   GLSL
	fragment GLSL
}

type GLSL struct {
	tokens Tokens
}

type Tokens []lexer.Token

func (s *Shader) String() string {
	return fmt.Sprintf("Shader(%d, vertex=%q, fragment=%q)", s.version, s.vertex.tokens, s.fragment.tokens)
}

func Build(tokenStream <-chan lexer.Token) string {
	tokens := make(Tokens, 0)

	var buffer bytes.Buffer

	shader := Shader{}

	for token := range tokenStream {
		if token.Typ == lexer.TokenVersionNumber {
			shader.version, _ = strconv.Atoi(token.Val)
		}
		tokens = append(tokens, token)
	}

	vertexStart := tokens.findFirst(lexer.TokenVertex) + 1
	vertexEnd := tokens.findFirst(lexer.TokenEnd)

	shader.vertex = GLSL{tokens[vertexStart:vertexEnd]}

	fragmentStart := tokens.findFirst(lexer.TokenFragment) + 1
	fragmentEnd := tokens.findLast(lexer.TokenEnd)

	shader.fragment = GLSL{tokens[fragmentStart:fragmentEnd]}

	for _, token := range tokens {
		if token.Typ == lexer.TokenVoid {
			buffer.WriteString("\n")
		} else if token.Typ == lexer.TokenGLSL {
			buffer.WriteString(token.Val)
		} else {
			buffer.WriteString(token.Typ.String())
			buffer.WriteString("(")
			buffer.WriteString(token.Val)
			buffer.WriteString(")")
		}
	}

	buffer.WriteString(shader.String())

	fmt.Println(shader)
	return shader.String()
}
