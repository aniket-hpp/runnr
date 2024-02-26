package config

var (
	ExtnDecl      []string                  //extension declarations
	VarsDecl      = make(map[string]string) //variable declarations
	CmdDecl       = make(map[string]string) //command declarations
	BuildDecl     []string                  //build declarations
	ModinFileData = make(map[string]string) //stores the file dir and file data
)
