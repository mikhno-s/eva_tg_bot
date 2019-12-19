package app

import (
	"log"
)

func checkErrorFatal(err error, cause string) {
	if err != nil {
		log.Fatalf("%s:\n%s", cause, err.Error())
	}
}
