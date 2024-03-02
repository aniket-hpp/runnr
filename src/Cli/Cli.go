package cli

import (
	"fmt"
	"io"
	"os"
	config "runnr/src/Config"
	internals "runnr/src/Internals"
	logerr "runnr/src/Logerr"
	"runtime"
	"strings"
)

var (
	Path            string
	MakeBuild       bool
	BuildArgs       string
	Args            []string
	FileName        string
	Extn            string
	Verbose         bool
	GetPreProcessed bool
	Slash           = "/"
	docsPath        string
	templatesPath   string
)

const (
	configFileName = "config.runnr"
	buildFileName  = "build.runnr"
	libraryPath    = ".runnr"
)

// Starts Processing Cli args
// Almost complete
func StartCli() {
	checkOS()
	// loops through args to check if single commands args are present or not
	i := 1
	for i < len(Args) {
		switch Args[i] {
		case "version": // prints current runnr version
			fmt.Println(internals.RUNNR_VERSION)
			os.Exit(0)

		case "help", "-h":
			{
				// if help has other args
				if i+1 < len(Args) {
					showHelp(Args[i+1])
				}

				// else grab home dir
				homeDir, _ := os.UserHomeDir()
				// readd default help.txt
				help, err := os.ReadFile(homeDir + Slash + docsPath + Slash + "help.txt")

				if err != nil {
					logerr.Log("docs path '%s' doesn't exists", homeDir+Slash+docsPath+Slash+"help.txt")
					os.Exit(1)
				}

				// print the data
				fmt.Println(string(help))
				os.Exit(0)
			}

		case "update":
			// TODO
			checkForUpdate()

		case "init", "-i":
			{
				// init has no args
				if i+1 >= len(Args) {
					logerr.Log("usage: runnr init <config/script>")
					os.Exit(1)
				}

				// else create Templates
				createTemplate(Args[i+1])
			}

		case "build", "-b":
			{
				// if buildd have args
				if i+1 < len(Args) {
					BuildArgs = Args[i+1]
				} else {
					BuildArgs = "runnr.null"
				}

				// set build to be true
				MakeBuild = true
				Verbose = true
				goto setPathandExit
			}

		case "verbose", "-v":
			Verbose = true
			i++

		case "preprocessor", "-p":
			GetPreProcessed = true
			i++

		case "config-path":
			fmt.Println(setPath())
			os.Exit(0)

		case "file", "-f":
			{
				if i+1 >= len(Args) {
					logerr.Log("no input filename for 'file'")
					os.Exit(1)
				}

				setFileName(i + 1)
				i += 2
			}

		default:
			{
				if i+1 < len(Args) {
					config.VarsDecl[Args[i]] = Args[i+1]
				} else {
					config.VarsDecl[Args[i]] = "runnr.null"
				}
				i += 2
			}
		}
	}

setPathandExit:
	Path = setPath()
}

// Function to set correct dir / or \\ based on os
func checkOS() {
	if runtime.GOOS == "windows" {
		Slash = "\\"
		fmt.Printf("%s\n", docsPath)
	}

	docsPath = libraryPath + Slash + "docs"
	templatesPath = libraryPath + Slash + "templates"
}

// Function to set config paths
// it check current path first else moves to default UserHomeDir
func setPath() string {
	// path selection
	path := configFileName
	cwd, _ := os.Getwd()

	// if user passed build through cli
	// then return build.runnr if exists in cwd
	if MakeBuild {
		if _, err := os.Stat(cwd + Slash + buildFileName); err != nil {
			logerr.Log("no 'build.runnr' file found in '%s' directory", cwd+Slash+buildFileName)
			os.Exit(1)
		}

		return cwd + Slash + buildFileName
	}

	// checks if a config file exists in current dir
	_, err := os.Stat(cwd + Slash + path)
	if err == nil {
		return cwd + Slash + path
	}

	// else default to UserHomeDir
	homeDir, _ := os.UserHomeDir()
	_, err = os.Stat(homeDir + Slash + path)
	// if config file doesn't exists
	if err != nil {
		os.Create(homeDir + Slash + path)
	}

	// return default path
	return homeDir + Slash + path
}

// Function to print sub docs
func showHelp(docName string) {
	// reads user home dir
	homeDir, _ := os.UserHomeDir()
	// tries to read the doc file
	help, err := os.ReadFile(homeDir + Slash + docsPath + Slash + docName + ".txt")

	// if fails
	if err != nil {
		logerr.Log("no docs found for '%s'", docName)
		os.Exit(1)
	}

	// else print it to the console
	fmt.Println(string(help))
	os.Exit(0)
}

func createTemplate(templateType string) {
	// get current working directory
	cwd, _ := os.Getwd()

	switch templateType {
	case "config", "script":
		{
			// by default chosing config file name
			fileName := configFileName

			// else switch to script file name
			if templateType == "script" {
				fileName = buildFileName
			}

			// check if a file already exists in cwd
			if _, err := os.Stat(cwd + Slash + fileName); err == nil {
				logerr.Warn("%s file already exists in '%s'", templateType, cwd+Slash+fileName)
				os.Exit(1)
			}

			// else grab the home directory
			homeDir, _ := os.UserHomeDir()
			// fetch the datas from template folder
			src, err := os.Open(homeDir + Slash + templatesPath + Slash + fileName)
			defer src.Close()

			if err != nil {
				logerr.Log("failed to open template file from '%s'", homeDir+Slash+templatesPath+Slash+fileName)
				os.Exit(1)
			}

			// create a file in cwd
			dest, err := os.Create(cwd + Slash + fileName)
			defer dest.Close()

			if err != nil {
				logerr.Log("failed to create %s file in '%s'", templateType, cwd+Slash+fileName)
				os.Exit(1)
			}

			// copy the contents of the file
			_, err = io.Copy(dest, src)

			if err != nil {
				logerr.Log("failed to copy data in file '%s'", cwd+Slash+fileName)
				os.Exit(1)
			}

			os.Exit(0)
		}

	default:
		// anything other than config/script is unknown
		logerr.Log("unknown option '%s' for 'init'", templateType)
		os.Exit(1)
	}
}

// Function to setup and store the cli args
func SetupArgs(args []string) {
	// if no cmd line args have been passed
	if len(args) <= 1 {
		logerr.Warn("no command-line argumensts\ntry: runnr help")
		os.Exit(1)
	}

	//else store all the args
	Args = args
}

// Function to check for update from an api
// TODO
func checkForUpdate() {
	logerr.Log("not implementedd yet")
	os.Exit(1)
}

// Function to parse the args
// It sets rest of the args in Vars Decl
func setFileName(i int) {
	arg := Args[i]
	// check if it is an file
	if index := strings.LastIndex(arg, "."); index != -1 {
		// check if it exists
		if _, err := os.Stat(arg); err != nil {
			logerr.Log("%s", err.Error())
			os.Exit(1)
		}

		// slicing and storing the value in Varibale Map
		FileName = arg[:index]
		config.VarsDecl["file"] = FileName
		Extn = arg[index+1:]
		config.VarsDecl["extn"] = Extn
	}
}
