package mantis_test

import (
	"testing"
	"container/list"
	"mantis"
)

func TestSaveList(t *testing.T) {
	l := list.New()

	var issue *mantis.Issue

	issue = new(mantis.Issue)
	issue.Id = 1
	issue.Level = "S"
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 2
	issue.Level = "A"
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 3
	issue.Level = "B"
	l.PushBack(issue)

	err := mantis.SaveList(l,"test.db")
	if(err != nil) {
		t.Fatal(err)
	}
}

func TestSaveListSameId(t *testing.T) {
	l := list.New()

	var issue *mantis.Issue

	issue = new(mantis.Issue)
	issue.Id = 1
	issue.Level = "S"
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 2
	issue.Level = "A"
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 3
	issue.Level = "B"
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 3
	issue.Level = "C"
	l.PushBack(issue)

	err := mantis.SaveList(l,"test1.db")
	if(err != nil) {
		t.Fatal(err)
	}
}
