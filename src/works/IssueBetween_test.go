package works_test

import (
	"testing"
	"time"
	"works"
)

func TestIssueBetween(t *testing.T) {
	works.IssueBetween()
}

func TestTime(t *testing.T)  {
	updated := "2018-05-18"
	lasttm, err := time.Parse("2006-01-02", updated)
	t.Logf("lasttm=[%s] err=%s\n" ,lasttm,err)

	t.Logf("updated=[%s] lasttm=[%s]\n" ,updated , lasttm.Format("2006-01-02 03:04:05 PM"))

}
