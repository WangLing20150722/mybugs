package main

import (
	"container/list"
	"flag"
	"log"
	"../mantis"
	"time"
	"../utils"
	"../works"
)

var (
	Version = ""
	Compile = ""
)

func main() {
	var start string
	flag.StringVar(&start, "start", "", "start time(included),eg:2018-4-1")

	var end string
	flag.StringVar(&end, "end", "", "end time(included),eg:2018-4-30")

	var listAction bool
	flag.BoolVar(&listAction, "list", true, "Update issue list")

	var pid string
	flag.StringVar(&pid, "pid", "0", "Set project id, default:0,all projects")

	var startPage int
	flag.IntVar(&startPage, "startpage", 0, "begin refresh from startpage, if it`s too big, some allIssues may lost but if too small, more times will be wasted. default:0")

	var detailAction bool
	flag.BoolVar(&detailAction, "detail", true, "Update issue details")

	var identifyAction bool
	flag.BoolVar(&identifyAction, "identify", true, "Identify issue owner")

	var failedAction bool
	flag.BoolVar(&failedAction,"failed",true,"Failed issue")

	var ownerFile string
	flag.StringVar(&ownerFile, "owner-config", "", "Json file path of onwer file")

	var outfile string
	flag.StringVar(&outfile, "output", "allIssues.xlsx", "Output excel file for identify")

	var username string
	flag.StringVar(&username, "username", "", "Username of bug system(mantis)")

	var password string
	flag.StringVar(&password, "password", "", "Password of bug system(mantis)")

	var forceRefresh bool
	flag.BoolVar(&forceRefresh, "force-refresh", false, "Force to refresh details")

	var version bool
	flag.BoolVar(&version, "version", false, "Show version")

	flag.Parse()

	if version {
		log.Printf("MyBus:\n	Version:%s\n	Compile:%s\n", Version, Compile)
		return
	}

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

	if listAction || detailAction {
		if username == "" || password == "" { //没有指定用户名密码，调试环境下不报错
			if utils.CONFIG.Username == "" || utils.CONFIG.Password == "" {
				log.Println("no username or password assign")
				flag.Usage()
				return
			}
		} else { //指定了用户名密码，优先使用命令行的
			utils.CONFIG.Username = username
			utils.CONFIG.Password = password
		}
	}

	var l *list.List
	if identifyAction {
		if ownerFile == "" {
			log.Println("Error -owner-config is must when 'identify'")
			flag.Usage()
			return
		} else {
			l = utils.ReadNameList(ownerFile)
			if l == nil {
				log.Println("Error -owner-config parse error, it should be a json string file")
				return
			}
		}
	}

	err = mantis.OpenDB()
	if err != nil {
		log.Printf("OpenDB error:%v\n", err)
		return
	}
	defer mantis.CloseDB()

	if listAction {
		err = works.RefreshListBetween(starttm, endtm, pid, startPage)
		if err != nil {
			log.Printf("List error:%v\n", err)
			return
		}
	}

	if detailAction {
		err = works.RefreshDetailsBetween(starttm, endtm, forceRefresh)
		if err != nil {
			log.Printf("Detail error:%v\n", err)
			return
		}
	}

	var allIssues *list.List
	if identifyAction && l != nil {
		issueList, err := works.IdentifyOwnerBetween(starttm, endtm, l)
		if err != nil {
			log.Printf("identify error:%v\n", err)
			return
		}

		//用于操作failed
		allIssues = list.New()
		allIssues.PushBackList(issueList)

		err = works.IssueOwnerList2Excel(issueList, outfile,"All_issues")
		if err != nil {
			log.Printf("identify write error:%v\n", err)
			return
		}
	}

	if failedAction && allIssues != nil{
		failedIssues := list.New()
		for e := allIssues.Front(); e != nil;e = e.Next()  {
			failedIssue := e.Value.(*works.IssueOwner)


			if failedIssue.Failed{
				failedIssues.PushBack(failedIssue)
			}
		}

		log.Println("failed allIssues length = ",failedIssues.Len())

		//works.IssueOwnerList2Excel(failedIssues,outfile,"failed_Issues")
		works.IssueOwnerFailedList2Excel(failedIssues,outfile,"failed_Issues")
		if err != nil {
			log.Printf("failedIssues write error:%v\n", err)
			return
		}

	}

	log.Println("Done!")
}
