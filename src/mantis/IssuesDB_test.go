package mantis_test

import (
	"testing"
	"container/list"
	"mantis"
	"time"
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

	err := mantis.SaveList(l)
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
	issue.FetchTime = time.Now()
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 2
	issue.Level = "A"
	issue.FetchTime = time.Now()
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 3
	issue.Level = "B"
	issue.FetchTime = time.Now()
	l.PushBack(issue)

	issue = new(mantis.Issue)
	issue.Id = 3
	issue.Level = "C"
	issue.FetchTime = time.Now()
	l.PushBack(issue)

	mantis.ClearList()
	err := mantis.SaveList(l)
	if(err != nil) {
		t.Fatal(err)
	}
}

func TestSaveDetail(t *testing.T) {

	detail := new(mantis.IssueDetail)
	detail.Id = 180263
	detail.FetchTime = time.Now()
	detail.History = `[{"DateModified":"2018-05-04T19:24:00Z","Username":"SQA_何纯贵_13421527372","Field":"NewIssue","Change":""},{"DateModified":"2018-05-04T19:24:00Z","Username":"SQA_何纯贵_13421527372","Field":"AssignedTo","Change":"=\u003e中间件_李辉_13923445825"},{"DateModified":"2018-05-04T19:24:00Z","Username":"SQA_何纯贵_13421527372","Field":"FileAdded:sitalog.rar","Change":""},{"DateModified":"2018-05-04T19:25:00Z","Username":"SQA_何纯贵_13421527372","Field":"FileAdded:Serial-COM4_0504_091352.log","Change":""}]`

	mantis.SaveDetail(detail)

	detail.Id = 180264
	detail.FetchTime = time.Now()
	mantis.SaveDetail(detail)

	detail.Id = 180263
	detail.FetchTime = time.Now()
	mantis.SaveDetail(detail)
}
