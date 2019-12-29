package app

import (
	"fmt"
	"runtime"
)

func checkError(err error, cause string) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		fmt.Printf("%s:\n%s:%d - %s\n", cause, file, line, err.Error())
	}

}
