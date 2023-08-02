package helper

import "log"

func Error(err error, msg string) {
	log.Println(msg)
	log.Fatal(err)
}
