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
	supply string
}

type Shader struct {
	Version  string
	vertex   Tokens
	fragment Tokens
	global   Tokens
	libs     []lib
	compiled glsl
}

type glsl struct {
	vertex   string
	fragment string
	provides []Tokens
	requests []Tokens
}

type lib struct {
	vertex   Tokens
	fragment Tokens
}

type varDef struct {
	typ  string
	name string
}

type libRes struct {
	dat Shader
	supply string
}

type Tokens []lexer.Token

const providePlaceholder = "___PROVIDE___REPLACE___HERE___\n"

func (s *Shader) String() string {
	return fmt.Sprintf("Shader(%s, vertex=%q, fragment=%q, global=%q)", s.Version, s.vertex, s.fragment, s.global)
}

func (s *Shader) GetVertex() string {
	return s.compiled.vertex
}

func (s *Shader) GetFragment() string {
	return s.compiled.fragment
}

func New(tokenStream <-chan lexer.Token) BuildStream {
	buildStream := BuildStream{
		make(chan error),
		make(chan LexRequest),
		make(chan Shader),
	}
	go build(tokenStream, buildStream)
	return buildStream
}

func build(tokenStream <-chan lexer.Token, buildStream BuildStream) {
	var shader *Shader
	var libs []lib

	libIndex := make(map[int]string)

	mainResponse := make(chan libRes)
	libResponse := make(chan libRes)
	reqPath := make(chan LexRequest)
	go generateShader(tokenStream, reqPath, mainResponse, buildStream.Err)

	libsPending := 0

loop:
	for {
		select {
		case s := <-mainResponse:
			shader = &s.dat
			if libsPending == 0 {
				break loop
			}
		case l := <-libResponse:
			libsPending -= 1
			libs = append(libs, lib{l.dat.vertex, l.dat.fragment})
			libIndex[len(libs)-1] = l.supply
			if libsPending == 0 && shader != nil {
				break loop
			}
		case dat := <-reqPath:
			libsPending += 1
			libStream := make(chan lexer.Token)
			dat.Answer = libStream
			buildStream.Request <- dat
			go generateShader(libStream, reqPath, libResponse, buildStream.Err)
		}
	}

	shader.injectLibs(libs)

	shader.buildVertex()
	shader.buildFragment()

	for _, request := range shader.compiled.requests {
		if !Contains(shader.compiled.provides, request) {
			buildStream.Err <- errors.New("Missing @provide statement for <" + request[0].Val + " " + request[1].Val + ">")
			return
		}
	}

	buildStream.Response <- *shader
}

func (shader *Shader) injectLibs(libs []lib) {
	shader.libs = libs
	injectLibVertex(shader)
	injectLibFragment(shader)
}

func (shader *Shader) buildVertex() {

	// vertex shader can only provide data to the fragment shader
	s, p, _ := buildGeneric(shader.vertex, shader.Version)
	shader.compiled.vertex = s
	shader.compiled.provides = p
}

func (shader *Shader) buildFragment() {

	// fragment shader can only request data from the vertex shader
	s, _, r := buildGeneric(shader.fragment, shader.Version)
	shader.compiled.fragment = s
	shader.compiled.requests = r
}
