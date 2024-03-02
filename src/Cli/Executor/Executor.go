package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	cli "runnr/src/Cli"
	config "runnr/src/Config"
	logerr "runnr/src/Logerr"
	"runtime"
)

// Function to execute given command string list
func Execute(cmds []string) {
	shell, ok := config.VarsDecl["shell"]

	if !ok {
		if runtime.GOOS == "windows" {
			shell = "powershell"
		} else {
			shell = "bash"
		}

		logerr.Warn("variable '$shell' is undefined in '.var' directive of file '%s', using default '%s'", cli.Path, shell)
	}

	// loops over the cmds
	for _, cmd := range cmds {
		var errBuff bytes.Buffer
		// if verbose then print the excuted cmds
		if cli.Verbose {
			fmt.Println(cmd)
		}

		// execute the command
		c := exec.Command(shell, "-c", cmd)
		c.Stderr = &errBuff
		res, err := c.Output()
		fmt.Print(string(res))

		// if err
		if err != nil {
			fmt.Printf("error message: %s", errBuff.String())
			fmt.Printf("exit status: %s\n", err.Error())
			os.Exit(1)
		}
	}
}
