// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
)

type Shader struct {
	version  int
	vertex   GLSL
	fragment GLSL
	global   GLSL
}

type GLSL struct {
	tokens Tokens
}

type Tokens []lexer.Token

func (s *Shader) String() string {
	return fmt.Sprintf("Shader(%d, vertex=%q, fragment=%q, global=%q)", s.version, s.vertex.tokens, s.fragment.tokens, s.global.tokens)
}

func Build(tokenStream <-chan lexer.Token) string {
	global := make(Tokens, 0)
	vertex := make(Tokens, 0)
	fragment := make(Tokens, 0)

	shader := Shader{}

	// phase 0 : global
	// phase 1 : vertex
	// phase 2 : fragment
	phase := 0
	for token := range tokenStream {
		switch token.Typ {
		case lexer.TokenVertex:
			phase = 1
		case lexer.TokenFragment:
			phase = 2
		case lexer.TokenEnd:
			phase = 0
		case lexer.TokenVoid:
			continue
		}

		switch phase {
		case 0:
			global = append(global, token)
		case 1:
			vertex = append(vertex, token)
		case 2:
			fragment = append(fragment, token)
		default:
			panic("unknow phase")
		}
	}

	shader.global = GLSL{global}
	shader.vertex = GLSL{vertex}
	shader.fragment = GLSL{fragment}

	return shader.String()
}
