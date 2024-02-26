package logerr

import (
	"fmt"
	lexer "runnr/src/Config/Lexer"
	tokens "runnr/src/Config/Tokens"
)

func Exepected(pos lexer.CursorPos, expected_tk int, got_tk int) bool {
	if expected_tk != got_tk {
		LogPos(pos)
		Log(fmt.Sprintf("expected '%s' but got '%s'\n", tokens.TokenNames[expected_tk], tokens.TokenNames[got_tk]))
		return false
	}

	return true
}
