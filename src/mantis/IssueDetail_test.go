package mantis_test

import (
	"testing"
	"mantis"
	"log"
)

func TestGetIssueDetail(t *testing.T) {
	_,err := mantis.Login("lihui02","asdfzxcv")
	log.Print("Login: ",err)

	_,err = mantis.SetProjectById("0")
	log.Print("SetProjectById: ",err)

	detail,err := mantis.GetIssueDetail(180263)
	if(err != nil) {
		log.Print(err)
	}

	log.Print(detail)
}
