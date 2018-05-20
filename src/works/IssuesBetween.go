package works

import (
	"time"
	"mantis"
	"log"
	"utils"
	"sync"
)

var DEBUG = true

func IssueListBetween(start,end time.Time,projectid string,startPage int,process int) error {
	var err error

	var startDay,endDay time.Time
	startDay,err = formatTime(start)
	if(err != nil) {
		log.Print("IssueListBetween formatTime failed:", err)
		return err
	}
	endDay,err = formatTime(end)
	if(err != nil) {
		log.Print("IssueListBetween formatTime failed:", err)
		return err
	}
	log.Printf("IssueListBetween %s --> %s of pid=%s startPage=%d\n",startDay.Format("2006-01-02"),endDay.Format("2006-01-02"),projectid,startPage)

	_,err = mantis.Login(utils.CONFIG.Username,utils.CONFIG.Password)
	if(err != nil) {
		log.Print("IssueListBetween Login failed:",err)
		return err
	} else {
		log.Println("IssueListBetween Login success")
	}

	_,err = mantis.SetProjectById(projectid)
	if(err != nil) {
		log.Print("IssueListBetween SetProjectById failed:",err)
		return err
	} else {
		log.Println("IssueListBetween SetProjectById success")
	}

	if(startPage < 0) {
		startPage = 0
	}
	page := startPage

	var mutex = &sync.Mutex{}
	var lastId int64
	finished := false

	var wg sync.WaitGroup

	downloadPageFunc := func (goid int) {
		for true {
			mutex.Lock()
			thispage := page
			page++
			mutex.Unlock()
			morepage, thisId, err := doDownloadPage(goid, thispage, false, endDay, startDay)

			if(err != nil) {
				continue
			}

			if (!morepage) {
				mutex.Lock()
				finished = true
				mutex.Unlock()
				break
			}

			mutex.Lock()
			if(thisId == lastId) {
				log.Printf("[%d]lastId =%d and curLastId=%d, may be the last page!\n",goid, lastId,thisId)
				finished = true
				mutex.Unlock()
				break
			} else {
				lastId = thisId
				mutex.Unlock()
			}
		}

		wg.Done()
	}

	if(process <= 0) {
		process = 1
	}

	//download 1st page to sort
	_, lastId, err = doDownloadPage(0,page,true,endDay,startDay)
	//TODO: add if more page here
	page++

	for i:=0;i<process;i++ {
		wg.Add(1)
		go downloadPageFunc(i)
	}

	wg.Wait()

	return nil
}

func doDownloadPage(goid int,page int,endDayWarning bool,endDay time.Time,startDay time.Time) (bool,int64,error) {
	l,err := mantis.ListBugs(page,"last_updated",false)
	log.Printf("[%d]IssueListBetween ListBugs page %d,err=%v",goid,page,err)
	mantis.SaveList(l)

	if ( DEBUG ) {
		front := l.Front().Value.(*mantis.Issue)
		back := l.Back().Value.(*mantis.Issue)
		log.Printf("	[%d]List Update: (%d)%s->(%d)%s\n",goid,front.Id,front.Updated,back.Id,back.Updated)
	}

	if(endDayWarning) {
		warning := true

		today,err := formatTime(time.Now())
		if(err == nil && today.Equal(endDay)) {
			//结束时间是今天，无需提示
			warning = false
		}

		if(warning) {
			endUpdated := l.Front().Value.(*mantis.Issue).Updated
			endtm, _ := time.Parse("2006-01-02", endUpdated)
			if(!endtm.After(endDay)) {
				log.Printf("[%d]You wants issues before %s but 1st got %s. Some issues may lost!\n",goid,endDay.Format("2006-01-02"),endUpdated)
			}
		}
	}

	lastUpdated := l.Back().Value.(*mantis.Issue).Updated
	lasttm, err := time.Parse("2006-01-02", lastUpdated)
	if(err != nil) {
		log.Printf("time.Parse %s error %v\n",lastUpdated,err)
		return false,0,err

	}

	morePage := true

	if(lasttm.Before(startDay)) {
		morePage =false;
	}

	thisId := l.Back().Value.(*mantis.Issue).Id

	return morePage,thisId,nil
}

/**
将精确的时间，转换成mantis List页面的时间显示:
2018-05-20 10:12:30  --> 2018-05-20
 */
func formatTime(exact time.Time) (time.Time,error) {
	return time.Parse("2006-01-02", exact.Format("2006-01-02"))
}
