package works

import (
	"log"
	"../mantis"
	"time"
	"../utils"
)

func RefreshDetailsBetween(start, end time.Time, forceRefresh bool) error {
	var err error

	var startDay, endDay time.Time
	startDay, err = utils.FormatTime2Day(start)
	if err != nil {
		log.Print("RefreshDetailsBwtween formatTime failed:", err)
		return err
	}
	endDay, err = utils.FormatTime2Day(end)
	if err != nil {
		log.Print("RefreshDetailsBwtween formatTime failed:", err)
		return err
	}
	log.Printf("RefreshDetailsBwtween %s --> %s\n", startDay.Format("2006-01-02"), endDay.Format("2006-01-02"))

	if !mantis.CheckLogin(nil) {
		_, err = mantis.Login(utils.CONFIG.Username, utils.CONFIG.Password)
		if err != nil {
			log.Printf("RefreshDetailsBetween Login failed %v\n", err)
			return err
		}
	}

	var issues []*mantis.Issue
	issues, err = mantis.EachList(startDay, endDay)
	if err != nil {
		log.Print("RefreshDetailsBwtween EachList failed:", err)
		return err
	}

	log.Printf("RefreshDetailsBwtween %d issues found in list:\n", len(issues))

	for i, issue := range issues {

		if forceRefresh {
			log.Printf("(%d/%d)Issue [%d] Updated at %s forceRefresh, do Refresh", i, len(issues), issue.Id, issue.Updated)
			doRefreshDetail(issue.Id)
		} else {
			var detail *mantis.IssueDetail
			detail, err = mantis.GetDetail(issue.Id)
			if err != nil {
				log.Printf("(%d/%d)Issue [%d] Updated at %s no details, do Refresh", i, len(issues), issue.Id, issue.Updated)
				doRefreshDetail(issue.Id)
			} else {
				ftday, _ := utils.FormatTime2Day(detail.FetchTime)
				updateTm, _ := time.Parse("2006-01-02", issue.Updated)

				if !ftday.After(updateTm) {
					log.Printf("(%d/%d)Issue [%d] Updated at %s Refreshed at %s, do Refresh\n", i, len(issues), issue.Id, issue.Updated, ftday.Format("2006-01-02"))
					doRefreshDetail(issue.Id)
				} else if detail.FetchTime.IsZero() || detail.History == "" || detail.Project == "" {
					log.Printf("(%d/%d)Issue [%d] detail not complted, do Refresh\n", i, len(issues), issue.Id)
					doRefreshDetail(issue.Id)
				} else {
					log.Printf("(%d/%d)Issue [%d] Updated at %s Refreshed at %s, will not Refresh\n", i, len(issues), issue.Id, issue.Updated, ftday.Format("2006-01-02"))
				}
			}
		}
	}

	return nil

}

const sleepTime,maxTryTime = 30,5
func doRefreshDetail(id int64) error {
	var err error
	var tryCount = 0



	LOOP:for i := 0; i < 3; i++ {
		var detail *mantis.IssueDetail
		detail, err = mantis.GetIssueDetail(id)

		if err != nil {
			log.Printf("doRefreshDetail failed:%v\n", err)
			continue
		}
		mantis.SaveDetail(detail)
		return nil
	}

	if tryCount<5{
		tryCount ++
		time.Sleep(time.Duration(tryCount*sleepTime)*time.Second)
		goto LOOP
	}

	return err
}

//func  doDetail(id int64)error  {
//
//}


