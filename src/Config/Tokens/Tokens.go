package tokens

// all available tokens
const (
	TK_DECLARE = iota
	TK_DEF     = iota
	TK_VAR     = iota
	TK_CMD     = iota
	TK_BUILD   = iota
	TK_PROJECT = iota
	TK_VERSION = iota
	TK_MODIN   = iota
	TK_JUMP    = iota

	TK_OPEN_CURLY_PAREN  = iota
	TK_CLOSE_CURLY_PAREN = iota

	TK_IDENTIFIER = iota
	TK_MODIFIER   = iota
	TK_ASSIGNMENT = iota
	TK_BACK_TICK  = iota
	TK_COMMA      = iota

	TK_EOF   = iota
	TK_ERROR = iota
	COUNT    = iota
)

var (
	TokenNames = [COUNT]string{
		"declaration directive",
		"definition directive",
		"variable directive",
		"command directive",
		"build directive",
		"project directive",
		"version directive",
		"module import directive",
		"preprocessor jump directive",

		"open curly-brace",
		"close curly-brace",

		"identifier name",
		"modifier name",
		"assignment",
		"backtick",
		"comma",

		"end-of-file",
		"unknown token",
	}
)
