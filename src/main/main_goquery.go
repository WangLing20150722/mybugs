package main

import (
	"flag"
	"log"
	"mantis"
	"time"
	"works"
)

func main() {
	var start string
	flag.StringVar(&start, "start", "", "start time(included),eg:2018-4-1")

	var end string
	flag.StringVar(&end, "end", "", "end time(included),eg:2018-4-30")

	var list bool
	flag.BoolVar(&list, "list", true, "Update issue list")

	var pid string
	flag.StringVar(&pid, "pid", "0", "Set project id, default:0,all projects")

	var startPage int
	flag.IntVar(&startPage, "startpage", 0, "begin refresh from startpage, if it`s too big, some issues may lost but if too small, more times will be wasted. default:0")

	var detail bool
	flag.BoolVar(&detail, "detail", true, "Update issue details")

	flag.Parse()

	var err error
	var starttm, endtm time.Time
	starttm, err = time.Parse("2006-01-02", start)
	if err != nil {
		log.Printf("Error -start %s\n", start)
		flag.Usage()
		return
	}

	endtm, err = time.Parse("2006-01-02", end)
	if err != nil {
		log.Printf("Error -end %s\n", end)
		flag.Usage()
		return
	}

	mantis.OpenDB()
	defer mantis.CloseDB()

	if list {
		err = works.RefreshListBetween(starttm, endtm, pid, startPage)
		if err != nil {
			log.Printf("List error:%v\n", err)
			return
		}
	}

	if detail {
		err = works.RefreshDetailsBetween(starttm, endtm)
		if err != nil {
			log.Printf("Detail error:%v\n", err)
			return
		}
	}
}
