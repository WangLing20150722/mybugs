package main

import (
	"log"
	"mantis"
	"container/list"
)


func main() {

	_,err := mantis.Login("lihui02","asdfzxcv")
	log.Print("Login: ",err)

	_,err = mantis.SetProjectById("0")
	log.Print("SetProjectById: ",err)

	var l *list.List
	l,err = mantis.ListBugs(0,2,"last_updated",false)
	log.Print("ListBugs: ",err)

	mantis.SaveList(l,"issues.db")

	mantis.DumpCookies()

}