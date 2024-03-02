package parser

import (
	"os"
	cli "runnr/src/Cli"
	config "runnr/src/Config"
	lexer "runnr/src/Config/Lexer"
	tokens "runnr/src/Config/Tokens"
	internals "runnr/src/Internals"
	logerr "runnr/src/Logerr"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// Main Parser struct to bind all related functions
type Parser struct {
	l lexer.Lexer
}

// Initializes the Parser & Lexer
// parameters: filename string, file-content string
func (p *Parser) Init(fileName, data string) {
	p.l.Init(fileName, data)
}

// Starts finding and parsing the data for the given extension
// parameter: extension *string
// returns: ExDef struct
func (p *Parser) ParseExtension(extn *string) map[string]string {
	// empty config file
	if len(p.l.Data) == 0 {
		logerr.Log("empty config file: '%s'", cli.Path)
		os.Exit(1)
	}

	// finding and parsing .declare
	if !p.findDeclaration() {
		logerr.Log("no '.declare' directive found in config file: '%s'", cli.Path)
		os.Exit(1)
	}

	// checking if extension exists
	haveExtn := slices.Index[[]string](config.ExtnDecl, *extn)
	if haveExtn == -1 {
		logerr.Log("declaration of '%s' not found in the config", *extn)
		os.Exit(1)
	}

	// reset the cursor
	p.l.ResetCursor()
	// return the parsed .def
	return *p.parseDefinition(extn)
}

// Function to find and parse .build directive
func (p *Parser) ParseBuild() {
	if p.findBuild() {
		p.parseBuild()
	}
}

// Find the ".declare" definition in the config file
// parameter: none
// returns: bool [true if found] else [false]
func (p *Parser) findDeclaration() bool {
	// flag for storing
	// if declare is found or not
	flag := false

	// loops thorugh entire file
	for !p.l.IsEnd() {
		switch p.l.GetNextToken() {
		// parse .declare if found
		case tokens.TK_DECLARE:
			flag = true
			p.parseDeclaration()

		// skip everything else
		case tokens.TK_VAR,
			tokens.TK_CMD,
			tokens.TK_BUILD:
			// skiping the entire {} block
			p.skipEntireBlock()

		// skips .def <extn>
		case tokens.TK_DEF:
			p.l.GetNextToken()

			if !logerr.Exepected(p.l.Pos, tokens.TK_MODIFIER, p.l.CurrentToken) {
				os.Exit(1)
			}

			p.skipEntireBlock()

		case tokens.TK_PROJECT:
			p.parseProject()

		case tokens.TK_VERSION:
			p.parseVersion()

		case tokens.TK_EOF:
			// end of file
			return flag

		// token jump <jump to content of another file>
		case tokens.TK_JUMP:
			// find and store the TK_DECLARE
			newFlag, _ := p.parseJump(tokens.TK_DECLARE)
			if !flag {
				flag = newFlag
			}

		default:
			{
				// any other token in global scope is nothing but error
				logerr.LogPos(p.l.CurrentTokenPos)
				logerr.Log("got unwanted declaration '%s' in global scope", p.l.GetCurrentRawLine())
				os.Exit(1)
			}
		}
	}

	return flag
}

// Parses ".declare" directives
// parameter: extension *string
func (p *Parser) parseDeclaration() {
	token := p.l.GetNextToken() // eat .declare

	// next token must be {
	if !logerr.Exepected(p.l.CurrentTokenPos, tokens.TK_OPEN_CURLY_PAREN, token) {
		os.Exit(1)
	}

	token = p.l.GetNextToken() //eat {

	// parsing tokens inside {}
	for !p.l.IsEnd() {
		// storing the identifier aka extension declaration
		if token == tokens.TK_IDENTIFIER {
			config.ExtnDecl = append(config.ExtnDecl, p.l.Identifier)
			token = p.l.GetNextToken()
		}

		// breaking the loop if it is an curly-brace }
		if token == tokens.TK_CLOSE_CURLY_PAREN {
			break
		}

		// any other token than comma or curly-braces
		if token != tokens.TK_COMMA {
			logerr.LogPos(p.l.CurrentTokenPos)
			logerr.Log("expected %s or %s but got %s",
				tokens.TokenNames[tokens.TK_COMMA], tokens.TokenNames[tokens.TK_CLOSE_CURLY_PAREN],
				tokens.TokenNames[token])
			os.Exit(1)
		}

		token = p.l.GetNextToken()
	}
}

// Find the ".def <exten>" definition the the config file
// parameter: extension *string
// returns: bool [true if found] else [false], extensionDefinition if found on another file, else nil
func (p *Parser) findDefinition(extn *string) (bool, *map[string]string) {
	for !p.l.IsEnd() {
		// checking tokens in file
		switch p.l.GetNextToken() {
		case tokens.TK_DEF:
			{
				p.l.GetNextToken() // eat .def

				if !logerr.Exepected(p.l.CurrentTokenPos, tokens.TK_MODIFIER, p.l.CurrentToken) {
					os.Exit(1)
				}

				// if identifier == extension, we got what we are looking for
				if p.l.Identifier == *extn {
					return true, nil
				}

				//else skip the entire {} block
				p.skipEntireBlock()
			}

		case tokens.TK_VAR:
			// parsing the variable declaration block
			p.parseVariableDeclarations()

		case tokens.TK_CMD, tokens.TK_BUILD, tokens.TK_DECLARE:
			// skiping the entire {} block
			p.skipEntireBlock()

		case tokens.TK_EOF:
			// end of file
			return false, nil

		case tokens.TK_PROJECT:
			p.parseProject()

		case tokens.TK_VERSION:
			p.parseVersion()

		case tokens.TK_JUMP:
			// .jump directive
			{
				// find and parse the TK_DEF in another file
				flag, extnDef := p.parseJump(tokens.TK_DEF)
				if flag {
					// if found return the definition
					return flag, extnDef
				}
			}

		default:
			{
				// any other token in global scope is nothing but error
				logerr.LogPos(p.l.CurrentTokenPos)
				logerr.Log("got unwanted declaration '%s' in global scope", p.l.GetCurrentRawLine())
				os.Exit(1)
			}
		}
	}

	// not found
	return false, nil
}

// Parses ".def" directives
// parameter: extension *string
// returns: *ExDef struct
func (p *Parser) parseDefinition(extn *string) *map[string]string {
	// finding the definition
	if found, def := p.findDefinition(extn); !found {
		logerr.Log("extension '%s' definition not found in config", *extn)
		os.Exit(1)
	} else {
		// if found check
		// if it exists on another file, i.e already parsed
		// then return the def
		if def != nil {
			return def
		}
	}

	// stores the data in an map
	extnDef := make(map[string]string)
	p.parseArgsVals(&extnDef)

	//retuns after parsing
	return &extnDef
}

// Function to skip entire {} block
func (p *Parser) skipEntireBlock() {
	// expect '{'
	if !logerr.Exepected(p.l.CurrentTokenPos, tokens.TK_OPEN_CURLY_PAREN, p.l.GetNextToken()) {
		os.Exit(1)
	}

	// skips all character until '}'
	for !p.l.IsEnd() && p.l.CurrentChar() != '}' {
		// eat '{'
		p.l.ForwardCursor()

		switch p.l.CurrentChar() {
		case '{':
			p.skipEntireBlock()

		case '`':
			p.l.ForwardCursor()
			p.parseBackTick()
		}
	}

	// eat `}`
	p.l.ForwardCursor()
}

// Function to parse content inside backticks `<data>...`
// returns: string [content inside the backtick]
func (p *Parser) parseBackTick() string {
	//stores actual data
	var data string

	//keep storing data until close '`' is found
	for p.l.CurrentChar() != '`' {
		data += string(p.l.CurrentChar())
		p.l.ForwardCursor()
	}

	p.l.ForwardCursor() // eat '`'
	return data
}

// Function to parse the value of modifier " = `<value>`"
// returns: string [content inside the backtick]
func (p *Parser) parseAssignment() string {
	// expects TK_ASSIGNMENT
	if !logerr.Exepected(p.l.Pos, tokens.TK_ASSIGNMENT, p.l.GetNextToken()) {
		os.Exit(1)
	}

	// eat =
	p.l.GetNextToken()

	// expects '`'
	if !logerr.Exepected(p.l.Pos, tokens.TK_BACK_TICK, p.l.CurrentToken) {
		os.Exit(1)
	}

	// if got '`' then pases its contains
	data := p.parseBackTick()

	//stores the final string with replaced data with variable
	//stores the index
	var finalData string
	i := 0

	//loops over the data
	for i < len(data) {
		//if the char is '$' {a variable}
		if data[i] == '$' {
			//strores the starting index
			i++
			beg := i

			//loops
			for i < len(data) && (unicode.IsLetter(rune(data[i])) || unicode.IsDigit(rune(data[i])) || data[i] == '_') {
				i++
			}

			//slice the actual variable from the string
			xvar := data[beg:i]

			//if the variable name starts with $env
			//its an env function call
			if xvar == "env" {
				//next char must be (
				if data[i] != '(' {
					logerr.Log("expected '(' after '$env' at line '%d: %s' in config file '%s'", cli.Path, p.l.Pos.Row, data)
					os.Exit(1)
				}

				i++ //eat (
				//store the begining
				envBeg := i

				//loop until )
				for i < len(data) && data[i] != ')' {
					i++
				}

				//retrieve the env and eat )
				finalData += os.Getenv(data[envBeg:i])
				i++
			} else {
				//checks and copies the variable value
				xval, contains := config.VarsDecl[xvar]

				//if variable is undefind
				if !contains {
					logerr.LogPos(p.l.Pos)
					logerr.Log("undefined variable '$%s'", xvar)
					os.Exit(1)
				}

				//appends the final data of the var
				finalData += xval
			}

			//if any space
			if i < len(data) {
				finalData += string(data[i])
			}
		} else {
			//just keep appending the string
			finalData += string(data[i])
		}

		i++
	}

	//returns the final data
	return finalData
}

// Parses ".var" directives and stores it in config.VarsDecl
func (p *Parser) parseVariableDeclarations() {
	p.parseArgsVals(&config.VarsDecl)
}

func (p *Parser) parseArgsVals(dataMap *map[string]string) {
	// eat .<directive>
	p.l.GetNextToken()

	// expects '{'
	if !logerr.Exepected(p.l.Pos, tokens.TK_OPEN_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}

	// eat '{'
	p.l.GetNextToken()

	for !p.l.IsEnd() {
		// expects TK_MODIFIER
		if !logerr.Exepected(p.l.Pos, tokens.TK_MODIFIER, p.l.CurrentToken) {
			os.Exit(1)
		}

		if _, ok := config.VarsDecl[p.l.Identifier]; !ok {
			// parse the TK_MODIFIER with its '= <data...'
			(*dataMap)[p.l.Identifier] = p.parseAssignment()
		} else {
			p.parseAssignment()
		}

		// eat TK_MODIFIER & its value
		p.l.GetNextToken()

		// breaks if it is '}'
		if p.l.CurrentToken == tokens.TK_CLOSE_CURLY_PAREN {
			break
		}

		// otherwise expects comma ','
		if !logerr.Exepected(p.l.Pos, tokens.TK_COMMA, p.l.CurrentToken) {
			os.Exit(1)
		}

		// skips comma
		p.l.GetNextToken()
	}
}

// Function to parse .jump directive
// parameters: token type to find and parse
// returns: bool if found something, and *map[string]string for the parse data
func (p *Parser) parseJump(token int) (bool, *map[string]string) {
	p.l.GetNextToken() // eat '.jump' directive

	// expects {
	if !logerr.Exepected(p.l.Pos, tokens.TK_OPEN_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}

	p.l.GetNextToken() // eat {

	// stores final result
	var flag bool
	var mainDef *map[string]string

	// loop through entire file
	for !p.l.IsEnd() {
		// if found backtick
		if p.l.CurrentToken == tokens.TK_BACK_TICK {
			// parse data inside backtick
			pathId := p.parseBackTick()
			var newP Parser

			// initialize a new parser for the data after checking
			if _, ok := config.ModinFileData[pathId]; !ok {
				logerr.Log("data not found for '%s' in file modin table.", pathId)
			}

			newP.Init(pathId, config.ModinFileData[pathId])

			// check which token to parse
			switch token {
			case tokens.TK_DECLARE:
				flag, mainDef = newP.findDeclaration(), nil

			case tokens.TK_DEF:
				// if def found in that file
				if found, def := newP.findDefinition(&cli.Extn); found {
					// with definition
					if def != nil {
						flag, mainDef = true, def
						break
					}

					// else parse it manually
					xExtnDef := make(map[string]string)
					newP.parseArgsVals(&xExtnDef)
					flag, mainDef = true, &xExtnDef
				}

			case tokens.TK_BUILD:
				newP.ParseBuild()
			}

			// move to next token
			p.l.GetNextToken()
		}

		// end of .jump
		if p.l.CurrentToken == tokens.TK_CLOSE_CURLY_PAREN {
			break
		}

		// else expects ,
		if !logerr.Exepected(p.l.Pos, tokens.TK_COMMA, p.l.CurrentToken) {
			os.Exit(1)
		}

		p.l.GetNextToken()
	}

	// finally return the data found
	return flag, mainDef
}

// Funtion to find .build directive while
// Parsing other directives
// Returns: true if a .build directive is found
func (p *Parser) findBuild() bool {
	// loops
	for !p.l.IsEnd() {
		switch p.l.GetNextToken() {
		case tokens.TK_CMD:
			// parses .cmd directive
			p.parseArgsVals(&config.CmdDecl)

		case tokens.TK_VAR:
			// parses .var directive
			p.parseVariableDeclarations()

		case tokens.TK_DECLARE, tokens.TK_DEF:
			// skiping .declare and .def
			{
				if p.l.CurrentToken == tokens.TK_DEF {
					p.l.GetNextToken()
					if !logerr.Exepected(p.l.Pos, tokens.TK_MODIFIER, p.l.CurrentToken) {
						os.Exit(1)
					}
				}

				p.skipEntireBlock()
			}

		case tokens.TK_JUMP:
			// parse .jump ~~ .modin directive
			p.parseJump(tokens.TK_BUILD)

		case tokens.TK_PROJECT:
			p.parseProject()

		case tokens.TK_VERSION:
			p.parseVersion()

		case tokens.TK_BUILD:
			// if found return true
			return true

		default:
			{
				// any other token in global scope is nothing but error
				logerr.LogPos(p.l.CurrentTokenPos)
				logerr.Log("got unwanted declaration '%s' in global scope", p.l.GetCurrentRawLine())
				os.Exit(1)
			}
		}
	}

	return false
}

// Function to parse .build directive
func (p *Parser) parseBuild() {
	p.l.GetNextToken() // eat .build directive

	// expects {
	if !logerr.Exepected(p.l.Pos, tokens.TK_OPEN_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}

	// loops
	for !p.l.IsEnd() {
		p.l.GetNextToken()

		// expects a modifier name
		if p.l.CurrentToken == tokens.TK_MODIFIER {
			config.BuildDecl = append(config.BuildDecl, p.l.Identifier)
			p.l.GetNextToken()
		}

		// exists if } end
		if p.l.CurrentToken == tokens.TK_CLOSE_CURLY_PAREN {
			break
		}

		// else must be ,
		if !logerr.Exepected(p.l.Pos, tokens.TK_COMMA, p.l.CurrentToken) {
			os.Exit(1)
		}
	}
}

// Function to parse .project directive
func (p *Parser) parseProject() {
	p.l.GetNextToken() // eat .project

	// expects {
	if !logerr.Exepected(p.l.Pos, tokens.TK_OPEN_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}

	p.l.GetNextToken() // eat {

	// expects a modifier name
	if !logerr.Exepected(p.l.Pos, tokens.TK_MODIFIER, p.l.CurrentToken) {
		os.Exit(1)
	}

	// stores it
	config.VarsDecl["project"] = p.l.Identifier

	// eats modifier
	p.l.GetNextToken()

	// expects }
	if !logerr.Exepected(p.l.Pos, tokens.TK_CLOSE_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}
}

// Function to parse .verion directive
func (p *Parser) parseVersion() {
	p.l.GetNextToken() // eat .version

	// expects {
	if !logerr.Exepected(p.l.Pos, tokens.TK_OPEN_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}

	p.l.GetNextToken() //eat {

	// expects vx_x_x
	if !logerr.Exepected(p.l.Pos, tokens.TK_MODIFIER, p.l.CurrentToken) {
		os.Exit(1)
	}

	// it must start with v and contains _
	if p.l.Identifier[0] == 'v' && strings.Contains(p.l.Identifier, "_") {
		// removing v
		p.l.Identifier = p.l.Identifier[1:]
		// spliting into token by _
		verList := strings.Split(p.l.Identifier, "_")

		// if less than 3 {maj, min, patch}
		if len(verList) != 3 {
			logerr.Log("wrong format for version in '.version' directive in '%s' file\nfromat should be: major_minor_patch", cli.Path)
			os.Exit(1)
		}

		// major version checking
		if v, err := strconv.Atoi(verList[0]); err != nil {
			logerr.Log("unexpected error: %s", err.Error())
			os.Exit(1)
		} else if v > internals.RUNNR_MAJ_VER {
			logerr.Log("script version of '%s' is not supported", cli.Path)
			os.Exit(1)
		}

		// minor version checking
		if v, err := strconv.Atoi(verList[1]); err != nil {
			logerr.Log("unexpected error: %s", err.Error())
			os.Exit(1)
		} else if v > internals.RUNNR_MIN_VER {
			logerr.Log("script version of '%s' is not supported", cli.Path)
			os.Exit(1)
		}
	} else {
		logerr.Log("wrong format for version in '.version' directive in '%s' file\nfromat should be: v<major>_<minor>_<patch>", cli.Path)
		os.Exit(1)
	}

	p.l.GetNextToken() //eat modifier

	if !logerr.Exepected(p.l.Pos, tokens.TK_CLOSE_CURLY_PAREN, p.l.CurrentToken) {
		os.Exit(1)
	}
}
