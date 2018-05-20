package works

import (
	"testing"
	"time"
)

func TestIssueBetween(t *testing.T) {
	today := time.Now()
	yesterday := today.AddDate(0,0,-1)
	err := IssueListBetween(yesterday,today,"0",-1)
	if(err != nil) {
		t.Fatal(err)
	}
}

func TestIssueBetweenOld(t *testing.T) {
	startday,_ := time.Parse("2006-01-02", "2018-04-01")
	endday,_ := time.Parse("2006-01-02", "2018-05-08")
	err := IssueListBetween(startday,endday,"565",-1)
	if(err != nil) {
		t.Fatal(err)
	}
}

func TestGenrateTime(t *testing.T)  {
	now := time.Now()
	today,err := formatTime(now)
	if(err != nil) {
		t.Fatal(err)
	}

	t.Logf("Today = %s",today.Format(time.RFC3339))

	if(today.Format("2006-01-02") != now.Format("2006-01-02")) {
		t.Fail()
	}

	if(today.Format("15:04:05Z07:00") != "00:00:00Z") {
		t.Fail()
	}

	yestodayNow := now.AddDate(0,0,-1)
	yestoday,err := formatTime(yestodayNow)
	if(err != nil) {
		t.Fatal(err)
	}

	t.Logf("Yestoday = %s",yestoday.Format(time.RFC3339))

	if(yestoday.Format("2006-01-02") != yestodayNow.Format("2006-01-02")) {
		t.Fail()
	}

	if(yestoday.Format("15:04:05Z07:00") != "00:00:00Z") {
		t.Fail()
	}
}
