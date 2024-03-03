package main

import (
	"fmt"
	"os"
	cli "runnr/src/Cli"
	cmdbuilder "runnr/src/Cli/CmdBuilder"
	executor "runnr/src/Cli/Executor"
	config "runnr/src/Config"
	parser "runnr/src/Config/Parser"
	preprocessor "runnr/src/Config/PreProcessor"
	logerr "runnr/src/Logerr"
)

func main() {
	// Parser struct
	var p parser.Parser

	// Stores cli args, sets config paths
	// & check inital single cli commands
	cli.SetupArgs(os.Args)
	cli.StartCli()

	// Reading Byte data from the Config File
	data, _ := os.ReadFile(cli.Path)
	processedData := preprocessor.PreProcess(cli.Path, string(data))

	if cli.GetPreProcessed {
		fmt.Println(processedData)
		os.Exit(0)
	}

	// Initialising the Parser
	p.Init(cli.Path, processedData)

	// if user passed "build" then we just execute commands in .build directive
	// and exit
	if cli.MakeBuild {
		p.ParseBuild()
		commands := cmdbuilder.GenerateBuildCmd(&config.BuildDecl, &config.CmdDecl)
		executor.Execute(commands)
		os.Exit(0)
	}

	// // if user inputed other args but not the input file
	if len(cli.Extn) == 0 {
		logerr.Log("no input files with extension")
		os.Exit(1)
	}

	// else we excute the commands realted to an file extension
	extnMap := p.ParseExtension(&cli.Extn)
	commands := cmdbuilder.GenerateExtnCmd(extnMap)
	executor.Execute(commands)
}
