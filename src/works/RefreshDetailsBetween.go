package works

import (
	"log"
	"mantis"
	"time"
	"utils"
)

func RefreshDetailsBetween(start, end time.Time) error {
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

	for _, issue := range issues {

		var detail *mantis.IssueDetail
		detail, err = mantis.GetDetail(issue.Id)
		if err != nil {
			log.Printf("Issue [%d] Updated at %s no details, do Refresh", issue.Id, issue.Updated)
			doRefreshDetail(issue.Id)
		} else {
			ftday, _ := utils.FormatTime2Day(detail.FetchTime)
			updateTm, _ := time.Parse("2006-01-02", issue.Updated)

			if !ftday.After(updateTm) {
				log.Printf("Issue [%d] Updated at %s Refreshed at %s, do Refresh\n", issue.Id, issue.Updated, ftday.Format("2006-01-02"))
				doRefreshDetail(issue.Id)
			} else {
				log.Printf("Issue [%d] Updated at %s Refreshed at %s, will not Refresh\n", issue.Id, issue.Updated, ftday.Format("2006-01-02"))
			}
		}
	}

	return nil

}

func doRefreshDetail(id int64) error {
	detail, err := mantis.GetIssueDetail(id)

	if err != nil {
		log.Printf("doRefreshDetail failed:%v\n", err)
		return err
	}
	mantis.SaveDetail(detail)
	return nil
}