package logerr

import (
	"fmt"
	"os"
)

func Warn(str string, args ...any) {
	fmt.Fprint(os.Stderr, "runnr: warn: ")
	fmt.Fprintf(os.Stderr, str+"\n", args...)
}

func Log(str string, args ...any) {
	fmt.Fprintf(os.Stderr, "runnr: error: "+str+"\n", args...)
}

func Abort(err error, str string, args ...any) {
	if err != nil {
		Log(str, args...)
		fmt.Fprintf(os.Stderr, "runnr: system error: %s", err.Error())
		os.Exit(1)
	}
}
