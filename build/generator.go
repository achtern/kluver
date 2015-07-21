// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"errors"
	"fmt"
	"github.com/achtern/kluver/lexer"
	"strings"
)

func generateShader(tokenStream chan lexer.Token, requests chan LexRequest, done chan libRes, err chan error, supply, filePath string) {
	shader := Shader{}

	shader.global = make(Tokens, 0)
	shader.vertex = make(Tokens, 0)
	shader.fragment = make(Tokens, 0)

	nextNameDecIsUse := false
	useSupply := ""

	// phase 0 : global
	// phase 1 : vertex
	// phase 2 : fragment
	phase := 0
	for token := range tokenStream {
		if token.Typ == lexer.TokenError {
			err <- errors.New(token.Val)
			return
		}
		if token.Typ == lexer.TokenExtendsName {
			fmt.Println("extending not yet supported.\nkill programm with ^C")
			requests <- LexRequest{token.Val, nil, supply}
		}

		if token.Typ == lexer.TokenUse {
			nextNameDecIsUse = true
		}

		if token.Typ == lexer.TokenNameDec && nextNameDecIsUse {
			useSupply = token.Val
		}

		if token.Typ == lexer.TokenImportPath {
			supply := ""

			if useSupply != "" {
				supply = useSupply
			}

			requests <- LexRequest{token.Val, nil, supply}
			useSupply = ""
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
			shader.Version = token.Val
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

	done <- libRes{shader, supply, filePath}
}

// buildGeneric compiles the generic shader, lib code has been injected already
func buildGeneric(tokens Tokens, version string) (string, []Tokens, []Tokens) {
	var sb StringBuffer
	sb.Append(buildHead(version))

	providePlaceholderInserted := false
	provides := make([]Tokens, 0)
	requests := make([]Tokens, 0)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		switch token.Typ {
		case lexer.TokenRequire:
			if !providePlaceholderInserted {
				// insert provides before the first uniform
				sb.Append(providePlaceholder)
				providePlaceholderInserted = true
			}
			sb.Append(generateRequire(tokens[i+1], tokens[i+2]))
			i += 2
		case lexer.TokenProvide:
			provides = append(provides, Tokens{tokens[i+1], tokens[i+2]})
			sb.Append(generateProvideSetting(tokens[i+1], tokens[i+2], tokens[i+3], tokens[i+4]))
			i += 4
		case lexer.TokenRequest:
			requests = append(requests, Tokens{tokens[i+1], tokens[i+2]})
			sb.Append(generateRequest(tokens[i+1], tokens[i+2]))
			i += 2
		case lexer.TokenYield:
			if tokens[i+1].Typ == lexer.TokenActionVar {
				i += 1
			}
			// the value has been replaced by the lib injector
			sb.Append(token.Val)
			i += 1
		case lexer.TokenWrite:
			tmpToken := lexer.Token{lexer.TokenVoid, 0, "fragColor" + tokens[i+2].Val}
			provides = append(provides, Tokens{tmpToken, tokens[i+4]})
			sb.Append(generateWriteAssignment("fragColor", tokens[i+2], tokens[i+4]))
			i += 4
		default:
			sb.Append(token.Val)
		}
	}

	compiled := sb.String()

	compiled = strings.Replace(compiled, providePlaceholder, generateProvideDecBlock(provides), -1)

	return compiled, provides, requests
}

func buildHead(version string) string {
	return "#version " + version + "\n"
}

func generateRequire(typ, name lexer.Token) string {
	return fmt.Sprintf("uniform %s %s", typ.Val, name.Val)
}

func generateProvideSetting(typ, name, assign, glslAction lexer.Token) string {
	return fmt.Sprintf("%s %s %s %s", typ.Val, name.Val, assign.Val, glslAction.Val)
}

func generateProvideDecBlock(provides []Tokens) string {
	var sb StringBuffer
	for _, tokens := range provides {
		sb.Append(generateProvideDec(tokens[0], tokens[1]))
	}

	return sb.String()
}

func generateProvideDec(typ, name lexer.Token) string {
	return fmt.Sprintf("out %s %s;\n", typ.Val, strings.Trim(name.Val, " "))
}

func generateRequest(typ, name lexer.Token) string {
	return fmt.Sprintf("in %s %s", typ.Val, name.Val)
}

func generateWriteAssignment(target string, slot, name lexer.Token) string {
	return fmt.Sprintf("%s%s = %s", target, slot.Val, name.Val)
}
