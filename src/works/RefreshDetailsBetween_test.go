package works_test

import (
	"testing"
	"time"
	"works"
	"mantis"
)

func TestRefreshDetailsBetweenOld(t *testing.T) {
	mantis.OpenDB()
	defer mantis.CloseDB()

	startday,_ := time.Parse("2006-01-02", "2018-04-01")
	endday,_ := time.Parse("2006-01-02", "2018-04-02")
	works.RefreshDetailsBetween(startday,endday)
}