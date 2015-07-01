// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package build

import (
	"fmt"
	"github.com/achtern/kluver/lexer"
)

func injectLibVertex(shader *Shader) {
outer:
	for {
		for i := 0; i < len(shader.vertex); i++ {
			if shader.vertex[i].Typ != lexer.TokenYield {
				if i == len(shader.vertex)-1 {
					break outer
				}
				continue
			}
			// contains everything before the yield token
			everythingBefore := append(Tokens(nil), shader.vertex[:i]...)

			// contains everything after the yield AND endstatement token
			everythingAfter := append(Tokens(nil), shader.vertex[i+2:]...)

			shader.vertex = everythingBefore

			include := false
			for _, lib := range shader.libs {
				for _, libToken := range lib.vertex {
					switch libToken.Typ {
					case lexer.TokenExport:
						include = true
						continue
					case lexer.TokenExportEnd:
						include = false
					}
					if include {
						shader.vertex = append(shader.vertex, libToken)
					}
				}
			}

			shader.vertex = append(shader.vertex, everythingAfter...)
			// break loop
			// this way we start at the beginning again
			break
		}
	}
}

func injectLibFragment(shader *Shader) {
	libGetterIdentifier := 0

	newFragment := make(Tokens, 0)

	include := false
	for _, lib := range shader.libs {
		for _, libToken := range lib.fragment {
			switch libToken.Typ {
			case lexer.TokenExport:
				include = true
				continue
			case lexer.TokenExportEnd:
				include = false
			case lexer.TokenGet:
				libToken.Val = fmt.Sprintf("vec4 get%d", libGetterIdentifier)
				libGetterIdentifier += 1
			}
			if include {
				newFragment = append(newFragment, libToken)
			}
		}
	}

	shader.fragment = append(newFragment, shader.fragment...)

	// for every yield token, we have to call the @get functions of all libs
	for i := 0; i < len(shader.fragment); i++ {
		if shader.fragment[i].Typ == lexer.TokenYield {
			var sb StringBuffer
			for x := range shader.libs {
				sb.append(fmt.Sprintf(
					"get%d(%s);", // fn call
					x, // libGetterIdentifier equiv.
					shader.fragment[i+1].Val)) // actionVar
			}
			shader.fragment[i].Val = sb.String()
		}
	}
}
