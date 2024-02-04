package utils

import "log"

func HandleError(e error) {
	if e != nil {
		log.Panicln(e)
	}
}
