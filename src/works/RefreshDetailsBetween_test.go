package works

import (
	"mantis"
	"testing"
	"time"
	"utils"
)

func TestRefreshDetailsBetweenOld(t *testing.T) {
	mantis.OpenDB()
	defer mantis.CloseDB()

	startday, _ := time.Parse("2006-01-02", "2018-04-01")
	endday, _ := time.Parse("2006-01-02", "2018-04-02")
	RefreshDetailsBetween(startday, endday, false)
}

func TestRefreshDetailsBetweenOldForce(t *testing.T) {
	mantis.OpenDB()
	defer mantis.CloseDB()

	startday, _ := time.Parse("2006-01-02", "2018-04-01")
	endday, _ := time.Parse("2006-01-02", "2018-04-02")
	RefreshDetailsBetween(startday, endday, true)
}

func TestRefreshDetailsOne(t *testing.T) {
	if !mantis.CheckLogin(nil) {
		_, err := mantis.Login(utils.CONFIG.Username, utils.CONFIG.Password)
		if err != nil {
			t.Logf("RefreshDetailsBetween Login failed %v\n", err)
			t.Fail()
		}
	}

	doRefreshDetail(174082)
	doRefreshDetail(170157)
	doRefreshDetail(164679)
	doRefreshDetail(170266)
	doRefreshDetail(166842)
	doRefreshDetail(173531)
	doRefreshDetail(174374)
	doRefreshDetail(165043)
}
