// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"fmt"
)

type token struct {
	Typ tokenType
	Pos int
	Val string
}

func (i token) String() string {
	switch i.Typ {
	case tokenEOF:
		return "tokenEOF"
	case tokenError:
		return i.Val
	}

	if len(i.Val) > 10 {
		return fmt.Sprintf("%s: %.10q[...]", i.Typ, i.Val)
	}
	return fmt.Sprintf("%s: %q", i.Typ, i.Val)
}

type tokenType int

const (
	tokenError tokenType = iota

	tokenEOF
	tokenVoid
	tokenVersion
	tokenVersionNumber
	tokenEndStatement
	tokenImport
	tokenImportPath
	tokenVertex
	tokenEnd
	tokenFragment
	tokenGLSL
	tokenYield
	tokenActionVar
	tokenProvide
	tokenRequire
	tokenRequest
	tokenAction
	tokenTypeDef
	tokenNameDec
	tokenAssign
	tokenGLSLAction
	tokenWrite
	tokenWriteOpenBracket
	tokenWriteCloseBracket
	tokenWriteSlot
)

func (i tokenType) String() string {
	switch i {
	case tokenVoid:
		return "tokenVoid"
	case tokenVersion:
		return "tokenVersion"
	case tokenVersionNumber:
		return "tokenVersionNumber"
	case tokenEndStatement:
		return "tokenEndStatement"
	case tokenImport:
		return "tokenImport"
	case tokenImportPath:
		return "tokenImport"
	case tokenVertex:
		return "tokenVertex"
	case tokenEnd:
		return "tokenEnd"
	case tokenFragment:
		return "tokenFragment"
	case tokenGLSL:
		return "tokenGLSL"
	case tokenAction:
		return "tokenAction"
	case tokenRequire:
		return "tokenRequire"
	case tokenProvide:
		return "tokenProvide"
	case tokenRequest:
		return "tokenRequest"
	case tokenYield:
		return "tokenYield"
	case tokenActionVar:
		return "tokenActionVar"
	case tokenTypeDef:
		return "tokenTypeDef"
	case tokenNameDec:
		return "tokenNameDec"
	case tokenAssign:
		return "tokenAssign"
	case tokenGLSLAction:
		return "tokenGLSLAction"
	case tokenWrite:
		return "tokenWrite"
	case tokenWriteOpenBracket:
		return "tokenWriteOpenBracket"
	case tokenWriteCloseBracket:
		return "tokenWriteCloseBracket"
	case tokenWriteSlot:
		return "tokenWriteSlot"
	default:
		return "unknown"
	}
}
