// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

func lexVoid(l *lexer) stateFn {
	for {
		if l.testPrefix(importLib, tokenVoid) {
			return lexImport
		}
		if l.testPrefix(vertex, tokenVoid) {
			return lexVertex
		}
		if l.testPrefix(fragment, tokenVoid) {
			return lexFragment
		}
		if l.next() == eof {
			break
		}
	}

	if l.pos > l.start {
		// reached end of file. Just void the stuff and exit
		l.emit(tokenVoid)
	}
	l.emit(tokenEOF)
	return nil
}

func lexEndStatementVoid(l *lexer) stateFn {
	return l.lexStatement(endStatement, tokenEndStatement, lexVoid)
}

func lexEndStatement(l *lexer) stateFn {
	return l.lexStatement(endStatement, tokenEndStatement, lexGLSL)
}

func lexImport(l *lexer) stateFn {
	return l.lexStatement(importLib, tokenImport, lexImportPath)
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
			if l.testPrefix(endStatement, tokenImportPath) {
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
	return l.lexStatement(vertex, tokenVertex, lexGLSL)
}

func lexFragment(l *lexer) stateFn {
	return l.lexStatement(fragment, tokenFragment, lexGLSL)
}

func lexEnd(l *lexer) stateFn {
	return l.lexStatement(end, tokenEnd, lexVoid)
}

func lexGLSL(l *lexer) stateFn {
	for {
		if l.testPrefix(action, tokenGLSL) {
			return lexAction
		}
		if l.testPrefix(end, tokenGLSL) {
			return lexEnd
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("unclosed GLSL block")
}

func lexAction(l *lexer) stateFn {
	l.lexStatement(action, tokenAction, nil)

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

	return l.errorf("unclosed action")
}

func lexRequire(l *lexer) stateFn {
	l.lexStatement(actionRequire, tokenRequire, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexRequest(l *lexer) stateFn {
	l.lexStatement(actionRequest, tokenRequest, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexProvide(l *lexer) stateFn {
	l.lexStatement(actionProvide, tokenProvide, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexTypeDef
}

func lexTypeDef(l *lexer) stateFn {
	for {
		if l.testPrefix(" ", tokenTypeDef) || l.testPrefix("\t", tokenTypeDef) {
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
		if l.testPrefix(endStatement, tokenNameDec) {
			return lexEndStatement
		}
		if l.testPrefix(actionAssign, tokenNameDec) {
			return lexActionAssign
		}
		if l.next() == eof {
			break
		}
		if l.hasPrefix(action) {
			break
		}
	}
	return l.errorf("incomplete name definition")
}

func lexActionAssign(l *lexer) stateFn {
	l.lexStatement(actionAssign, tokenAssign, nil)
	if isSpace(l.next()) {
		l.ignore()
	}
	return lexGLSLAction
}

func lexGLSLAction(l *lexer) stateFn {
	for {
		if l.testPrefix(endStatement, tokenGLSLAction) {
			return lexEndStatement
		}
		if l.next() == eof {
			break
		}
	}
	return l.errorf("incomplete glsl action assignment")
}

func lexYield(l *lexer) stateFn {
	return l.lexStatement(actionYield, tokenYield, lexActionVar)
}

func lexActionVar(l *lexer) stateFn {
	if l.hasPrefix(endStatement) {
		return lexEndStatement
	}
	for {
		if isSpace(l.next()) {
			l.ignore()
		}
		if l.testPrefix(endStatement, tokenActionVar) {
			return lexEndStatement
		}
	}
	return lexEndStatement
}

func lexWrite(l *lexer) stateFn {
	return l.lexStatement(actionWrite, tokenWrite, lexWriteOpenBracket)
	return lexWriteOpenBracket
}

func lexWriteOpenBracket(l *lexer) stateFn {
	if l.hasPrefix(writeOpenBracket) {
		return l.lexStatement(writeOpenBracket, tokenWriteOpenBracket, lexWriteSlot)
	}

	return l.errorf("<%s> expected after write action", writeOpenBracket)
}

func lexWriteCloseBracket(l *lexer) stateFn {
	if l.hasPrefix(writeCloseBracket) {
		return l.lexStatement(writeCloseBracket, tokenWriteCloseBracket, lexActionVar)
	}

	return l.errorf("<%s> expected after write slot", writeCloseBracket)
}

func lexWriteSlot(l *lexer) stateFn {
	if !l.isNumber() {
		return l.errorf("bad number syntax for write slot: %q", l.input[l.start:l.pos])
	}
	l.emit(tokenWriteSlot)
	return lexWriteCloseBracket
}
