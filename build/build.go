// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
	"github.com/achtern/kluver/lexer"
)

type BuildStream struct {
	Err      chan error
	Request  chan LexRequest
	Response chan Shader
}

type LexRequest struct {
	Path   string
	Answer chan lexer.Token
}

type Shader struct {
	version  string
	vertex   Tokens
	fragment Tokens
	global   Tokens
	compiled GLSL
}

type GLSL struct {
	vertex   string
	fragment string
	provides []Tokens
	requests []Tokens
}

type varDef struct {
	typ  string
	name string
}

type Tokens []lexer.Token

const providePlaceholder = "___PROVIDE___REPLACE___HERE___\n"

func (s *Shader) String() string {
	return fmt.Sprintf("Shader(%s, vertex=%q, fragment=%q, global=%q)", s.version, s.vertex, s.fragment, s.global)
}

func (s *Shader) GetVertex() string {
	return s.compiled.vertex
}

func (s *Shader) GetFragment() string {
	return s.compiled.fragment
}

func New(tokenStream <-chan lexer.Token) BuildStream {
	buildStream := BuildStream{
		nil,
		make(chan LexRequest),
		nil,
	}
	go build(tokenStream, buildStream)
	return buildStream
}

func build(tokenStream <-chan lexer.Token, buildStream BuildStream) {
	shader := Shader{}

	shader.global = make(Tokens, 0)
	shader.vertex = make(Tokens, 0)
	shader.fragment = make(Tokens, 0)
	

	// phase 0 : global
	// phase 1 : vertex
	// phase 2 : fragment
	phase := 0
	for token := range tokenStream {
		if token.Typ == lexer.TokenError {
			buildStream.Err <- errors.New(token.Val)
			return
		}

		if token.Typ == lexer.TokenImportPath {
			libStream := make(chan lexer.Token)
			buildStream.Request <- LexRequest{token.Val,libStream}
		}

		switch token.Typ {
		case lexer.TokenVertex:
			phase = 1
			continue
		case lexer.TokenFragment:
			phase = 2
			continue
		case lexer.TokenEnd:
			phase = 0
			continue
		case lexer.TokenVoid:
			continue
		}

		if token.Typ == lexer.TokenVersionNumber {
			shader.version = token.Val
			continue
		}

		if token.Typ == lexer.TokenAction {
			// we do need the action tokens after lexing
			continue
		}

		switch phase {
		case 0:
			shader.global = append(shader.global, token)
		case 1:
			shader.vertex = append(shader.vertex, token)
		case 2:
			shader.fragment = append(shader.fragment, token)
		default:
			panic("unknow phase")
		}
	}

	// shader.global = global
	// shader.vertex = vertex
	// shader.fragment = fragment

	shader.buildVertex()
	shader.buildFragment()

	for _, request := range shader.compiled.requests {
		if !contains(shader.compiled.provides, request) {
			buildStream.Err <- errors.New("Missing @provide statement for <" + request[0].Val + " " + request[1].Val + ">")
			return
		}
	}

	fmt.Println(shader.GetVertex())
}

func (shader *Shader) buildVertex() {

	// vertex shader can only provide data to the fragment shader
	s, p, _ := buildGeneric(shader.vertex, shader.version)
	shader.compiled.vertex = s
	shader.compiled.provides = p
}

func (shader *Shader) buildFragment() {

	// fragment shader can only request data from the vertex shader
	s, _, r := buildGeneric(shader.fragment, shader.version)
	shader.compiled.fragment = s
	shader.compiled.requests = r
}
