// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

import (
	"fmt"
)

type TokenType int

type Token struct {
	Typ TokenType
	Pos int
	Val string
}

func (i Token) String() string {
	switch i.Typ {
	case TokenEOF:
		return "TokenEOF"
	case TokenError:
		return i.Val
	}

	if len(i.Val) > 10 {
		return fmt.Sprintf("%s: %.10q[...]", i.Typ, i.Val)
	}
	return fmt.Sprintf("%s: %q", i.Typ, i.Val)
}

func (i Token) equals(other Token) bool {
	return i.compatible(other) && i.Pos == other.Pos
}

func (i Token) compatible(other Token) bool {
	return i.Typ == other.Typ && i.Val == other.Val
}

const (
	TokenError TokenType = iota

	TokenEOF
	TokenVoid
	TokenVersion
	TokenVersionNumber
	TokenExtends
	TokenExtendsName
	TokenEndStatement
	TokenImport
	TokenImportPath
	TokenUse
	TokenExport
	TokenExportEnd
	TokenVertex
	TokenEnd
	TokenFragment
	TokenGLSL
	TokenYield
	TokenActionVar
	TokenProvide
	TokenRequire
	TokenRequest
	TokenAction
	TokenTypeDef
	TokenNameDec
	TokenAssign
	TokenGLSLAction
	TokenWrite
	TokenActionOpenBracket
	TokenActionCloseBracket
	TokenWriteSlot
	TokenGet
)

func (i TokenType) String() string {
	switch i {
	case TokenVoid:
		return "TokenVoid"
	case TokenVersion:
		return "TokenVersion"
	case TokenVersionNumber:
		return "TokenVersionNumber"
	case TokenExtends:
		return "TokenExtends"
	case TokenExtendsName:
		return "TokenExtendsName"
	case TokenEndStatement:
		return "TokenEndStatement"
	case TokenImport:
		return "TokenImport"
	case TokenImportPath:
		return "TokenImportPath"
	case TokenUse:
		return "TokenUse"
	case TokenExport:
		return "TokenExport"
	case TokenExportEnd:
		return "TokenExportEnd"
	case TokenVertex:
		return "TokenVertex"
	case TokenEnd:
		return "TokenEnd"
	case TokenFragment:
		return "TokenFragment"
	case TokenGLSL:
		return "TokenGLSL"
	case TokenAction:
		return "TokenAction"
	case TokenRequire:
		return "TokenRequire"
	case TokenProvide:
		return "TokenProvide"
	case TokenRequest:
		return "TokenRequest"
	case TokenYield:
		return "TokenYield"
	case TokenActionVar:
		return "TokenActionVar"
	case TokenTypeDef:
		return "TokenTypeDef"
	case TokenNameDec:
		return "TokenNameDec"
	case TokenAssign:
		return "TokenAssign"
	case TokenGLSLAction:
		return "TokenGLSLAction"
	case TokenWrite:
		return "TokenWrite"
	case TokenActionOpenBracket:
		return "TokenActionOpenBracket"
	case TokenActionCloseBracket:
		return "TokenActionCloseBracket"
	case TokenWriteSlot:
		return "TokenWriteSlot"
	case TokenGet:
		return "TokenGet"
	default:
		return "unknown"
	}
}
