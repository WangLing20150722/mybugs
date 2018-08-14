package mantis_test

import (
	"log"
	"../mantis"
	"testing"
	"../utils"
)

func TestListBugsOrder(t *testing.T) {

	_, err := mantis.Login(utils.CONFIG.Username, utils.CONFIG.Password)
	log.Print("Login: ", err)

	_, err = mantis.SetProjectById("0")
	log.Print("SetProjectById: ", err)

	l, err := mantis.ListBugs(199, "last_updated", false)

	e := l.Back()
	log.Printf("%d - %s\n", e.Value.(*mantis.Issue).Id, e.Value.(*mantis.Issue).Updated)

	e = l.Front()
	log.Printf("%d - %s\n", e.Value.(*mantis.Issue).Id, e.Value.(*mantis.Issue).Updated)
}
