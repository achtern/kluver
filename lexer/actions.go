// Copyright 2015 Christian Gärtner. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package lexer

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

	if l.hasPrefix(actionExportEnd) {
		return lexExportEnd
	}

	if l.hasPrefix(actionExport) {
		return lexExport
	}

	if l.hasPrefix(actionGet) {
		return lexGet
	}

	if l.hasPrefix(actionTemplateEnd) {
		return lexTemplateEnd
	}

	if l.hasPrefix(actionTemplate) {
		return lexTemplate
	}

	if l.hasPrefix(actionSupplyEnd) {
		return lexSupplyEnd
	}

	if l.hasPrefix(actionSupply) {
		return lexSupply
	}

	return l.errorf("unclosed action")
}
