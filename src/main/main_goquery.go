package main

import (
	"mantis"
	"log"
)


func main() {
	_,err := mantis.Login("lihui02","asdfzxcv")
	log.Print(err)

	_,err = mantis.SetProjectById("0")
	log.Print(err)

	mantis.DumpCookies()



}