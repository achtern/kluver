// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
	"strings"
)

func generateRequire(action, typ, name lexer.Token) string {
	return fmt.Sprintf("uniform %s %s", typ.Val, name.Val)
}

func generateProvideSetting(typ, name, assign, glslAction lexer.Token) string {
	return fmt.Sprintf("%s %s %s %s", typ.Val, name.Val, assign.Val, glslAction.Val)
}

func generateProvideDecBlock(provides []Tokens) string {
	var sb StringBuffer
	for _, tokens := range provides {
		sb.append(generateProvideDec(tokens[0], tokens[1]))
	}

	return sb.String()
}

func generateProvideDec(typ, name lexer.Token) string {
	return fmt.Sprintf("out %s %s;\n", typ.Val, strings.Trim(name.Val, " "))
}
