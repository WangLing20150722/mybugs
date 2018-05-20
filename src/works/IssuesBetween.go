package works

import (
	"log"
	"mantis"
	"time"
	"utils"
)

var DEBUG = true

func IssueListBetween(start, end time.Time, projectid string, startPage int) error {
	var err error

	var startDay, endDay time.Time
	startDay, err = utils.FormatTime2Day(start)
	if err != nil {
		log.Print("IssueListBetween formatTime failed:", err)
		return err
	}
	endDay, err = utils.FormatTime2Day(end)
	if err != nil {
		log.Print("IssueListBetween formatTime failed:", err)
		return err
	}
	log.Printf("IssueListBetween %s --> %s of pid=%s startPage=%d\n", startDay.Format("2006-01-02"), endDay.Format("2006-01-02"), projectid, startPage)

	_, err = mantis.Login(utils.CONFIG.Username, utils.CONFIG.Password)
	if err != nil {
		log.Print("IssueListBetween Login failed:", err)
		return err
	} else {
		log.Println("IssueListBetween Login success")
	}

	_, err = mantis.SetProjectById(projectid)
	if err != nil {
		log.Print("IssueListBetween SetProjectById failed:", err)
		return err
	} else {
		log.Println("IssueListBetween SetProjectById success")
	}

	if startPage < 0 {
		startPage = 0
	}
	page := startPage

	var lastId int64

	for true {
		l, err := mantis.ListBugs(page, "last_updated", false)
		log.Printf("IssueListBetween ListBugs page %d,err=%v", page, err)
		mantis.SaveList(l)

		if DEBUG {
			front := l.Front().Value.(*mantis.Issue)
			back := l.Back().Value.(*mantis.Issue)
			log.Printf("	List Update: (%d)%s->(%d)%s\n", front.Id, front.Updated, back.Id, back.Updated)
		}

		//first page
		//page != 0是需要的，但无意义因为不存在，不加这个条件为方便构造测试用例
		if page == startPage /*&& page != 0 */ {
			warning := true

			today, err := utils.FormatTime2Day(time.Now())
			if err == nil && today.Equal(endDay) {
				//结束时间是今天，无需提示
				warning = false
			}

			if warning {
				endUpdated := l.Front().Value.(*mantis.Issue).Updated
				endtm, _ := time.Parse("2006-01-02", endUpdated)
				if !endtm.After(endDay) {
					log.Printf("You wants issues before %s but 1st got %s. Some issues may lost!\n", endDay.Format("2006-01-02"), endUpdated)
				}
			}
		}

		lastUpdated := l.Back().Value.(*mantis.Issue).Updated
		lasttm, err := time.Parse("2006-01-02", lastUpdated)
		if err != nil {
			return err

		}

		if lasttm.Before(startDay) {
			break
		} else {
			page++
		}

		curLastId := l.Back().Value.(*mantis.Issue).Id
		if curLastId == lastId {
			log.Printf("lastId =%d and curLastId=%d, may be the last page!\n", lastId, curLastId)
			break
		} else {
			lastId = curLastId
		}
	}

	return nil
}
