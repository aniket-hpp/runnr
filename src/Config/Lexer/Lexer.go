package lexer

import (
	tokens "runnr/src/Config/Tokens"
	"unicode"
)

// shows current cursor position
type CursorPos struct {
	FileName string
	Cursor   int
	Row      int
	LineBeg  int
}

// main lexer struct
type Lexer struct {
	Data            string    // file data
	Len             int       // file length
	CurrentToken    int       // current parsed token
	CurrentTokenPos CursorPos //current token position
	Identifier      string    // current identifier
	Pos             CursorPos // current position in file
}

// initializes the lexer struct
func (l *Lexer) Init(FileName, Data string) {
	l.Data = Data
	l.Pos.FileName = FileName
	l.Pos.Cursor = 0
	l.Pos.Row = 0
	l.Pos.LineBeg = 0
	l.Len = len(l.Data)

	l.CurrentTokenPos.Cursor = 0
	l.CurrentTokenPos.LineBeg = 0
	l.CurrentTokenPos.Row = 0
	l.CurrentTokenPos.FileName = FileName
}

// returns current character under cursor
func (l *Lexer) CurrentChar() byte {
	if l.IsEnd() {
		return 0
	}

	return l.Data[l.Pos.Cursor]
}

// moves the cursor forward across the file data [one character at a time]
func (l *Lexer) ForwardCursor() {
	if !l.IsEnd() {
		l.Pos.Cursor++

		// updates the row pointer and line beg pointer
		if !l.IsEnd() && l.CurrentChar() == '\n' {
			l.Pos.Row++
			l.Pos.LineBeg = l.Pos.Cursor + 1
		}
	}
}

// resets the cursor ro the begining
func (l *Lexer) ResetCursor() {
	l.Init(l.Pos.FileName, l.Data)
}

func (l *Lexer) SetTokenPos() {
	l.CurrentTokenPos = CursorPos{
		Cursor:   l.Pos.Cursor,
		Row:      l.Pos.Row,
		LineBeg:  l.Pos.LineBeg,
		FileName: l.Pos.FileName,
	}
}

// removes spaces or unneccessary character from left
func (l *Lexer) trimLeft() {
	for !l.IsEnd() && unicode.IsSpace(rune(l.CurrentChar())) {
		l.ForwardCursor()
	}
}

// drops a line [for comments]
func (l *Lexer) dropLine() {
	for !l.IsEnd() && l.CurrentChar() != '\n' {
		l.ForwardCursor()
	}
}

// returns if the cursor is at the end of the file or not
func (l *Lexer) IsEnd() bool {
	return l.Pos.Cursor >= l.Len
}

// generates tokens and return its token type and stores its values in identifier
func (l *Lexer) getToken() int {
	// always remove unwanted character
	l.trimLeft()
	l.SetTokenPos()

	// if the line is a comment
	if l.CurrentChar() == '#' {
		for l.CurrentChar() == '#' {
			l.dropLine()
			l.trimLeft()
		}
	}

	// then checking if there is data
	if l.IsEnd() {
		return tokens.TK_EOF
	}

	// if a token starts with alpha character or underscore
	if unicode.IsLetter(rune(l.CurrentChar())) || l.CurrentChar() == '_' {
		beg := l.Pos.Cursor

		for !l.IsEnd() && (unicode.IsLetter(rune(l.CurrentChar())) ||
			unicode.IsDigit(rune(l.CurrentChar())) ||
			l.CurrentChar() == '_') {
			l.ForwardCursor()
		}

		l.Identifier = l.Data[beg:l.Pos.Cursor]
		return tokens.TK_MODIFIER
	}

	// otherwise switches to different char
	switch l.CurrentChar() {
	case '.': // special directives or extension
		{
			l.ForwardCursor()
			beg := l.Pos.Cursor

			for !l.IsEnd() && unicode.IsLetter(rune(l.CurrentChar())) {
				l.ForwardCursor()
			}

			directive := l.Data[beg:l.Pos.Cursor]

			switch directive {
			case "modin":
				return tokens.TK_MODIN

			case "declare":
				return tokens.TK_DECLARE

			case "jump":
				return tokens.TK_JUMP

			case "var":
				return tokens.TK_VAR

			case "cmd":
				return tokens.TK_CMD

			case "def":
				return tokens.TK_DEF

			case "build":
				return tokens.TK_BUILD

			case "project":
				return tokens.TK_PROJECT

			case "version":
				return tokens.TK_VERSION

			default:
				l.Identifier = directive
				return tokens.TK_IDENTIFIER
			}
		}

	// special characters
	case '{':
		l.ForwardCursor()
		return tokens.TK_OPEN_CURLY_PAREN

	case '}':
		l.ForwardCursor()
		return tokens.TK_CLOSE_CURLY_PAREN

	case '=':
		l.ForwardCursor()
		return tokens.TK_ASSIGNMENT

	case ',':
		l.ForwardCursor()
		return tokens.TK_COMMA

	case '`':
		l.ForwardCursor()
		return tokens.TK_BACK_TICK
	}

	// else return error
	return tokens.TK_ERROR
}

// stores and returns the next token
func (l *Lexer) GetNextToken() int {
	l.CurrentToken = l.getToken()
	return l.CurrentToken
}

func (l *Lexer) GetCurrentRawLine() string {
	beg := l.Pos.Cursor
	l.dropLine()
	return l.Data[beg:l.Pos.Cursor]
}
