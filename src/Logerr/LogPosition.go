package logerr

import (
	"fmt"
	lexer "runnr/src/Config/Lexer"
)

func LogPos(pos lexer.CursorPos) {
	fmt.Printf("%s:%d:%d ", pos.FileName, pos.Row+1, pos.Cursor-pos.LineBeg+1)
}
