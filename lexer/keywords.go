package lexer

const (
	eof           = -1
	endStatement  = ";"
	version       = "#version"
	importLib     = "@import"
	vertex        = "#---VERTEX---#"
	fragment      = "#---FRAGMENT---#"
	end           = "#---END---#"
	action        = "@"
	actionRequire = "require"
	actionProvide = "provide"
	actionYield   = "yield"
	actionRequest = "request"
	actionWrite   = "write"
	actionExport  = "export"
	actionGet     = "get"
	actionAssign  = "="

	writeOpenBracket  = "("
	writeCloseBracket = ")"
	exportBlockOpen   = "{"
	exportBlockClose  = "}"
)
