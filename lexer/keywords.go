package lexer

const (
	eof             = -1
	endStatement    = ";"
	version         = "#version"
	extends         = "#extends"
	importLib       = "@import"
	useLib          = "@use"
	useFrom         = "from"
	vertex          = "#---VERTEX---#"
	fragment        = "#---FRAGMENT---#"
	end             = "#---END---#"
	action          = "@"
	actionRequire   = "require"
	actionProvide   = "provide"
	actionYield     = "yield"
	actionRequest   = "request"
	actionWrite     = "write"
	actionExport    = "export"
	actionExportEnd = "exportEnd"
	actionGet       = "get"
	actionAssign    = "="

	actionOpenBracket  = "("
	actionCloseBracket = ")"
)
