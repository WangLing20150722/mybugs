package works

import (
	"time"
	"mantis"
	"log"
)

func IssueBetween() error {

	_,err := mantis.Login("lihui02","asdfzxcv")
	log.Print("Login: ",err)

	_,err = mantis.SetProjectById("0")
	log.Print("SetProjectById: ",err)

	starttm,_ := time.Parse("2006-01-02", "2018-05-18")
	//endtm,_ := time.Parse("2006-01-02", "2018-05-17")

	page := 0
	for true {
		l,err := mantis.ListBugs(page,"last_updated",false)
		log.Print("ListBugs: ",err)
		mantis.SaveList(l)

		if ( true ) {
			front := l.Front().Value.(*mantis.Issue)
			back := l.Back().Value.(*mantis.Issue)
			log.Printf("id:%d->%d update %s->%s\n",front.Id,back.Id, front.Updated,back.Updated)
		}

		updated := l.Back().Value.(*mantis.Issue).Updated
		lasttm, err := time.Parse("2006-01-02", updated)
		if(err != nil) {
			return err

		}

		log.Printf("Last Update=%s\n" ,lasttm.Format("2006-01-02 03:04:05 PM"))

		if(lasttm.Before(starttm)) {
			break;
		} else {
			page++
		}
	}

	return nil
}
