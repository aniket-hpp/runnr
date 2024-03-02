package cmdbuilder

import (
	"os"
	cli "runnr/src/Cli"
	config "runnr/src/Config"
	logerr "runnr/src/Logerr"
)

// Function to build entier command/s for the
// given extension map
func GenerateExtnCmd(extnMap map[string]string) []string {
	// stores the commands in a list
	var cmdList []string
	var cmd1, cmd2 string
	// ordered array for checking the keys
	propList := [5]string{"exec", "flagsb", "out", "file", "flagsa"}

	// loops through the keys
	for i := 0; i < len(propList); i++ {
		// adding the filename w extension in cmd string
		if propList[i] == "file" {
			cmd1 += config.VarsDecl["file"] + "." + config.VarsDecl["extn"] + " "
			continue
		}

		// rest of the value of cmd
		if val, ok := extnMap[propList[i]]; ok {
			cmd1 += val + " "
		}
	}

	// finally appends it to the list
	cmdList = append(cmdList, cmd1)

	// for run and args command
	if val, ok := extnMap["run"]; ok && val == "y" {
		cmd2 += "." + cli.Slash + config.VarsDecl["file"]

		// wraps the args value in quote
		if val, ok = config.VarsDecl["args"]; ok {
			cmd2 += " \"" + val + "\""
		}

		// again appends it to list
		cmdList = append(cmdList, cmd2)
	}

	// return the list
	return cmdList
}

// Function to generate commands from .build & .cmd directives
func GenerateBuildCmd(buildList *[]string, cmdList *map[string]string) []string {
	// main string list
	var cmds []string

	// if no args for build
	if cli.BuildArgs == "runnr.null" {
		// loop over list in .build
		for _, build := range *buildList {
			// checking if the cmd exists in .cmd
			if val, ok := (*cmdList)[build]; ok {
				cmds = append(cmds, val)
			} else {
				logerr.Log("undefined command '%s' in '.cmd' directive in script file '%s'", build, cli.Path)
				os.Exit(1)
			}
		}
	} else {
		// if user have specied a particular cmd to execute
		// checking if it exists in .cmd
		if val, ok := (*cmdList)[cli.BuildArgs]; ok {
			cmds = append(cmds, val)
		} else {
			logerr.Log("undefined command '%s' in '.cmd' directive in script file '%s'", cli.BuildArgs, cli.Path)
		}
	}

	// finally return the cmd list
	return cmds
}
