// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

func lexVoid(l *lexer) stateFn {
	for {
		if l.testPrefix(version, TokenVoid) {
			return lexVersion
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

func lexAction(l *lexer) stateFn {
	l.lexStatement(action, TokenAction, nil)

	if l.hasPrefix(actionRequire) {
		return lexRequire
	}

	if l.hasPrefix(actionProvide) {
		return lexProvide
	}

	if l.hasPrefix(actionYield) {
		return lexYield
	}

	if l.hasPrefix(actionRequest) {
		return lexRequest
	}

	if l.hasPrefix(actionWrite) {
		return lexWrite
	}

	if l.hasPrefix(actionExport) {
		return lexExport
	}

	return l.errorf("unclosed action")
}

func lexRequire(l *lexer) stateFn {
	l.lexStatement(actionRequire, TokenRequire, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexRequest(l *lexer) stateFn {
	l.lexStatement(actionRequest, TokenRequest, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexProvide(l *lexer) stateFn {
	l.lexStatement(actionProvide, TokenProvide, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexExport(l *lexer) stateFn {
	l.lexStatement(actionExport, TokenExport, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexExportBlockOpen
}

func lexExportBlockOpen(l *lexer) stateFn {
	if l.hasPrefix(exportBlockOpen) {
		return l.lexStatement(exportBlockOpen, TokenExportBlockOpen, getLexGLSLBlock(TokenExportGLSL, lexExportBlockClose, exportBlockClose, "incomplete export block", true))
	}

	return l.errorf("<%s> expected after export action", exportBlockOpen)
}

func lexExportBlockClose(l *lexer) stateFn {
	if l.hasPrefix(exportBlockClose) {
		return l.lexStatement(exportBlockClose, TokenExportBlockClose, lexGLSL)
	}

	return l.errorf("<%s> expected after GLSL export", exportBlockClose)

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
	if isSpace(l.next()) {
		l.ignore()
	}
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
	if isSpace(l.next()) {
		l.ignore()
	}
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

func getLexGLSLBlock(token TokenType, next stateFn, blockTerminator, errorMsg string, allowLineBreaks bool) stateFn {
	return func (l *lexer) stateFn {
		for {
			if l.testPrefix(blockTerminator, token) {
				return next
			}
			if !allowLineBreaks && l.next() == '\n' {
				break
			}
			if l.next() == eof {
				break
			}
		}
		return l.errorf(errorMsg)
	}
}

func lexYield(l *lexer) stateFn {
	return l.lexStatement(actionYield, TokenYield, lexActionVar)
}

func lexActionVar(l *lexer) stateFn {
	if l.hasPrefix(endStatement) {
		return lexEndStatement
	}
	for {
		if isSpace(l.next()) {
			l.ignore()
		}
		if l.testPrefix(endStatement, TokenActionVar) {
			return lexEndStatement
		}
	}
	return lexEndStatement
}

func lexWrite(l *lexer) stateFn {
	return l.lexStatement(actionWrite, TokenWrite, lexWriteOpenBracket)
	return lexWriteOpenBracket
}

func lexWriteOpenBracket(l *lexer) stateFn {
	if l.hasPrefix(writeOpenBracket) {
		return l.lexStatement(writeOpenBracket, TokenWriteOpenBracket, lexWriteSlot)
	}

	return l.errorf("<%s> expected after write action", writeOpenBracket)
}

func lexWriteCloseBracket(l *lexer) stateFn {
	if l.hasPrefix(writeCloseBracket) {
		return l.lexStatement(writeCloseBracket, TokenWriteCloseBracket, lexActionVar)
	}

	return l.errorf("<%s> expected after write slot", writeCloseBracket)
}

func lexWriteSlot(l *lexer) stateFn {
	if !l.isNumber() {
		return l.errorf("bad number syntax for write slot: %q", l.input[l.start:l.pos])
	}
	l.emit(TokenWriteSlot)
	return lexWriteCloseBracket
}
