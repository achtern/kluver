// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

func lexVoid(l *lexer) stateFn {
	for {
		if l.testPrefix(version, TokenVoid) {
			return lexVersion
		}
		if l.testPrefix(extends, TokenVoid) {
			return lexExtends
		}
		if l.testPrefix(importLib, TokenVoid) {
			return lexImport
		}
		if l.testPrefix(vertex, TokenVoid) {
			return lexVertex
		}
		if l.testPrefix(fragment, TokenVoid) {
			return lexFragment
		}
		if l.testPrefix(useLib, TokenVoid) {
			return lexUse
		}
		if l.next() == eof {
			break
		}
	}

	if l.pos > l.start {
		// reached end of file. Just void the stuff and exit
		l.emit(TokenVoid)
	}
	l.emit(TokenEOF)
	return nil
}

func lexVersion(l *lexer) stateFn {
	l.lexStatement(version, TokenVersion, nil)
	for {
		if isSpace(l.next()) {
			l.ignore()
		} else {
			break
		}
	}
	return lexVersionNumber
}

func lexVersionNumber(l *lexer) stateFn {
	if !l.isNumber() {
		return l.errorf("bad number syntax for version number: %q", l.input[l.start:l.pos])
	}
	l.emit(TokenVersionNumber)
	return lexVoid
}

func lexExtends(l *lexer) stateFn {
	return l.lexStatement(extends, TokenExtends, lexExtendsName)
}

func lexExtendsName(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isSpace(r):
			l.ignore()
		case r == eof:
			return l.errorf("unclosed #extends statement")
		case r == '\n':
			l.backup() // we do not want to have the line break in the name
			l.emit(TokenExtendsName)
			l.next()   // advance forward
			l.ignore() // and ignore the line break
			return lexVoid
		}
	}
}

func lexEndStatementVoid(l *lexer) stateFn {
	return l.lexStatement(endStatement, TokenEndStatement, lexVoid)
}

func lexEndStatement(l *lexer) stateFn {
	return l.lexStatement(endStatement, TokenEndStatement, lexGLSL)
}

func lexImport(l *lexer) stateFn {
	return l.lexStatement(importLib, TokenImport, lexImportPath)
}

func lexImportPath(l *lexer) stateFn {
	for {
		if l.hasPrefix(endStatement) {
			return lexEndStatementVoid
		}

		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed import statement")
		case isSpace(r):
			l.ignore()
		}

		for {
			if l.testPrefix(endStatement, TokenImportPath) {
				return lexEndStatementVoid
			}
			if l.next() == eof {
				break
			}
		}
		return l.errorf("unclosed import statement")
	}
}

func lexUse(l *lexer) stateFn {
	return l.lexStatement(useLib, TokenUse, lexUseName)
}

func lexUseName(l *lexer) stateFn {
	l.ignoreSpace()
	for {
		if l.testPrefix(useFrom, TokenNameDec) {
			return lexUseFrom
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("unclosed use statement")
}

func lexUseFrom(l *lexer) stateFn {
	return l.lexStatement(useFrom, TokenUseFrom, lexImportPath)
}

func lexVertex(l *lexer) stateFn {
	return l.lexStatement(vertex, TokenVertex, lexGLSL)
}

func lexFragment(l *lexer) stateFn {
	return l.lexStatement(fragment, TokenFragment, lexGLSL)
}

func lexEnd(l *lexer) stateFn {
	return l.lexStatement(end, TokenEnd, lexVoid)
}

func lexGLSL(l *lexer) stateFn {
	for {
		if l.testPrefix(action, TokenGLSL) {
			return lexAction
		}
		if l.testPrefix(end, TokenGLSL) {
			return lexEnd
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("unclosed GLSL block")
}

func lexRequire(l *lexer) stateFn {
	l.lexStatement(actionRequire, TokenRequire, nil)
	l.ignoreSpace()
	return lexTypeDef
}

func lexRequest(l *lexer) stateFn {
	l.lexStatement(actionRequest, TokenRequest, nil)
	l.ignoreSpace()
	return lexTypeDef
}

func lexProvide(l *lexer) stateFn {
	l.lexStatement(actionProvide, TokenProvide, nil)
	l.ignoreSpace()
	return lexTypeDef
}

func lexExport(l *lexer) stateFn {
	l.lexStatement(actionExport, TokenExport, nil)
	l.ignoreSpace()
	return lexGLSL
}

func lexExportEnd(l *lexer) stateFn {
	l.lexStatement(actionExportEnd, TokenExportEnd, nil)
	l.ignoreSpace()
	return lexGLSL
}

func lexTemplate(l *lexer) stateFn {
	return l.lexStatement(actionTemplate, TokenTemplate, lexTemplateName)
}

func lexTemplateName(l *lexer) stateFn {
	l.ignoreSpace()
	for {
		if l.testPrefix(pointer, TokenNameDec) {
			return getLexPointer(lexGLSL)
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("unclosed template statement")
}

func lexTemplateEnd(l *lexer) stateFn {
	return l.lexStatement(actionTemplateEnd, TokenTemplateEnd, lexGLSL)
}

func lexSupply(l *lexer) stateFn {
	return l.lexStatement(actionSupply, TokenSupply, lexSupplyName)
}

func lexSupplyName(l *lexer) stateFn {
	l.ignoreSpace()
	for {
		if l.testPrefix(pointer, TokenNameDec) {
			return getLexPointer(lexGLSL)
		}
		if l.testPrefix(colon, TokenNameDec) {
			return getLexColon(lexSupplyName) // lexSupplyExtends would be the same method
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("unclosed supply statement")
}

func lexSupplyEnd(l *lexer) stateFn {
	return l.lexStatement(actionSupplyEnd, TokenSupplyEnd, lexGLSL)
}

func lexGet(l *lexer) stateFn {
	return l.lexStatement(
		actionGet,
		TokenGet,
		getLexActionOpenBracket(
			getLexActionVar(
				actionCloseBracket,
				getLexActionCloseBracket(
					lexGLSL,
					"<%> expected after get variable")),
			"<%> expected after get action"))
}

func lexTypeDef(l *lexer) stateFn {
	for {
		if l.testPrefix(" ", TokenTypeDef) || l.testPrefix("\t", TokenTypeDef) {
			return lexNameDec
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("incomplete type definition")
}

func lexNameDec(l *lexer) stateFn {
	l.ignoreSpace()
	for {
		if l.testPrefix(endStatement, TokenNameDec) {
			return lexEndStatement
		}
		if l.testPrefix(actionAssign, TokenNameDec) {
			return lexActionAssign
		}

		if l.next() == eof {
			break
		}
		if l.hasPrefix(action) {
			break
		}
		if l.peek() == '\n' {
			break
		}
	}
	return l.errorf("incomplete name definition")
}

func lexActionAssign(l *lexer) stateFn {
	l.lexStatement(actionAssign, TokenAssign, nil)
	l.ignoreSpace()
	return lexGLSLAction
}

func lexGLSLAction(l *lexer) stateFn {
	for {
		if l.testPrefix(endStatement, TokenGLSLAction) {
			return lexEndStatement
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("incomplete glsl action assignment")
}

func lexYield(l *lexer) stateFn {
	return l.lexStatement(actionYield, TokenYield, getLexActionVar(endStatement, lexEndStatement))
}

func lexWrite(l *lexer) stateFn {
	return l.lexStatement(actionWrite, TokenWrite, getLexActionOpenBracket(lexWriteSlot, "<%s> expected after write action"))
}

func lexWriteSlot(l *lexer) stateFn {
	if !l.isNumber() {
		return l.errorf("bad number syntax for write slot: %q", l.input[l.start:l.pos])
	}
	l.emit(TokenWriteSlot)
	return getLexActionCloseBracket(getLexActionVar(endStatement, lexEndStatement), "<%s> expected after write slot action")
}

func getLexGLSLBlock(token TokenType, next stateFn, blockTerminator, errorMsg string, allowLineBreaks bool) stateFn {
	return func(l *lexer) stateFn {
		for {
			if l.testPrefix(blockTerminator, token) {
				return next
			}
			if !allowLineBreaks && l.peek() == '\n' {
				l.next()
				break
			}
			if l.next() == eof {
				break
			}
		}
		return l.errorf(errorMsg)
	}
}

func getLexActionOpenBracket(next stateFn, errorMsg string) stateFn {
	return func(l *lexer) stateFn {
		if l.hasPrefix(actionOpenBracket) {
			return l.lexStatement(actionOpenBracket, TokenActionOpenBracket, next)
		}
		return l.errorf(errorMsg, actionOpenBracket)
	}
}

func getLexActionCloseBracket(next stateFn, errorMsg string) stateFn {
	return func(l *lexer) stateFn {
		if l.hasPrefix(actionCloseBracket) {
			return l.lexStatement(actionCloseBracket, TokenActionCloseBracket, next)
		}

		return l.errorf(errorMsg, actionCloseBracket)
	}
}

func getLexActionVar(terminator string, next stateFn) stateFn {
	return func(l *lexer) stateFn {
		if l.hasPrefix(terminator) {
			return next
		}
		for {
			if l.testPrefix(terminator, TokenActionVar) {
				return next
			}
			if isSpace(l.next()) {
				l.ignore()
			}
		}
		return next
	}
}

func getLexPointer(next stateFn) stateFn {
	return func(l *lexer) stateFn {
		l.ignoreSpace()
		return l.lexStatement(pointer, TokenPointer, next)
	}
}

func getLexColon(next stateFn) stateFn {
	return func(l *lexer) stateFn {
		l.ignoreSpace()
		return l.lexStatement(colon, TokenColon, next)
	}
}
