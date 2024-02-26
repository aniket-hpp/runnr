package preprocessor

import (
	"os"
	cli "runnr/src/Cli"
	config "runnr/src/Config"
	lexer "runnr/src/Config/Lexer"
	tokens "runnr/src/Config/Tokens"
	logerr "runnr/src/Logerr"
)

// func to pre-process module imports in the read config file
// parameter: data string
// returns string [processed config file]
func PreProcess(data string) string {
	// if file is empty
	if len(data) == 0 {
		return data
	}

	// initializing an lexer for the data
	var l lexer.Lexer
	l.Init(cli.Path, data)

	// stored processed Data
	var processedData string
	beg := 0

	// loop through file
	for !l.IsEnd() {
		// saves intial position
		tokenBeg := l.Pos.Cursor
		// store the row of the token before .modin
		rowBeg := l.Pos.Row

		//if the token is MODIN
		if l.GetNextToken() == tokens.TK_MODIN {
			// store read data until now
			processedData += data[beg:tokenBeg]
			// store current row of the .modin directive
			rowEnd := l.Pos.Row
			// adds the necessary lines in between
			for i := 0; i < rowEnd-rowBeg; i++ {
				processedData += "\n"
			}

			// again store the row before parsing inside the block
			rowBeg = l.Pos.Row
			// copy the .jump tag while replacing .modin{} block
			jumpPaths := parseModin(&l)
			// stores the row after parsing the block
			rowEnd = l.Pos.Row

			// if len of jumpPaths is 0, it is unecessary to add the block
			if len(jumpPaths) > 0 {
				processedData += "\n.jump {\n"

				for i := 0; i < len(jumpPaths); i++ {
					processedData += "\t`" + jumpPaths[i] + "`"

					if i+1 < len(jumpPaths) {
						processedData += ",\n"
					}

				}

				processedData += "\n}"
			} else {
				// if the .modin was empty
				// then the file is replace with approprite number of
				// next line

				// why we need this??
				// to thorugh error in the lines where it needs.
				// if we mismatch or removes lines then position to thorugh the line error
				// will change too
				for i := 0; i < rowEnd-rowBeg; i++ {
					processedData += "\n"
				}
			}

			// reset the cursor position to current char i.e, after .modin{} block
			beg = l.Pos.Cursor
		}

		if l.CurrentToken == tokens.TK_JUMP {
			logerr.LogPos(l.Pos)
			logerr.Log("internal token '.jump' is not allowed to declared in the runnr file")
			os.Exit(1)
		}

		//if we get BACK_TICK we skips it while ignoring its content
		if l.CurrentToken == tokens.TK_BACK_TICK {
			for !l.IsEnd() && l.CurrentChar() != '`' {
				l.ForwardCursor()
			}

			// eat '`'
			l.ForwardCursor()
		}
	}

	//return processedData with leftovers
	return processedData + data[beg:]
}

// func to Parse Modin block .modin{}
// parameter *lexer.Lexer
// returns string [parsed data of .modin{} files]
func parseModin(l *lexer.Lexer) []string {

	// stores current .modin position
	tkPos := l.CurrentTokenPos

	// expects {
	if !logerr.Exepected(l.Pos, tokens.TK_OPEN_CURLY_PAREN, l.GetNextToken()) {
		os.Exit(1)
	}

	var filePaths []string
	// loop thorugh contents
	for !l.IsEnd() {
		// eat {
		token := l.GetNextToken()

		// if ` then it contains a path
		if token == tokens.TK_BACK_TICK {
			// store the postion after `
			beg := l.Pos.Cursor
			// loop until get closing `
			for !l.IsEnd() && l.CurrentChar() != '`' {
				l.ForwardCursor()
			}

			//check if the path is of iteself or root
			if l.Data[beg:l.Pos.Cursor] == cli.Path {
				logerr.Log("import cycle is not allowed in '%s' runnr file", cli.Path)
				os.Exit(1)
			}

			// slice the path from it
			filePaths = append(filePaths, l.Data[beg:l.Pos.Cursor])
			l.ForwardCursor()        //eat '`'
			token = l.GetNextToken() // move to next token
		}

		// if token is } then we have reached .modin end
		if token == tokens.TK_CLOSE_CURLY_PAREN {
			break
		}

		// else it must be a comma
		if !logerr.Exepected(l.Pos, tokens.TK_COMMA, l.CurrentToken) {
			os.Exit(1)
		}
	}

	// loops though the path we got
	for _, path := range filePaths {
		fileData, err := os.ReadFile(path)
		// if file does not exists
		if err != nil {
			logerr.LogPos(tkPos)
			logerr.Log("failed to read file '%s' in runnr file: '%s'", path, cli.Path)
			os.Exit(1)
		}

		// else we will preprocess the file data we just read
		// yup this is a recursion process, ez
		config.ModinFileData[path] = PreProcess(string(fileData))
	}

	// finally return the file
	return filePaths
}
