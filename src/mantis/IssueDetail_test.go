package mantis_test

import (
	"log"
	"mantis"
	"testing"
	"utils"
)

func TestGetIssueDetail(t *testing.T) {
	_, err := mantis.Login(utils.CONFIG.Username, utils.CONFIG.Password)
	log.Print("Login: ", err)

	_, err = mantis.SetProjectById("0")
	log.Print("SetProjectById: ", err)

	detail, err := mantis.GetIssueDetail(180263)
	if err != nil {
		log.Print(err)
	}

	log.Print(detail)
}
