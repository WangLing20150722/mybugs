package works

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"mantis"
	"regexp"
	"strings"
	"time"
	"utils"
)

type IssueOwner struct {
	Id              int64
	Project         string
	Level           string
	Summary         string
	Status          string
	LastModify      string
	LastModifyTime  time.Time
	LastAssignOutTo string
	LastFix         string
	InTime          time.Time //首次分入时间
	OutTime         time.Time //分出或者解决时间
}

func IdentifyOwnerBetween(start, end time.Time, shortOnwers *list.List) (*list.List, error) {
	var err error

	var startDay, endDay time.Time
	startDay, err = utils.FormatTime2Day(start)
	if err != nil {
		log.Print("IdentifyOwnerBetween formatTime failed:", err)
		return nil, err
	}
	endDay, err = utils.FormatTime2Day(end)
	if err != nil {
		log.Print("IdentifyOwnerBetween formatTime failed:", err)
		return nil, err
	}
	log.Printf("IdentifyOwnerBetween %s --> %s\n", startDay.Format("2006-01-02"), endDay.Format("2006-01-02"))

	var issues []*mantis.Issue
	issues, err = mantis.EachList(startDay, endDay)
	if err != nil {
		log.Print("IdentifyOwnerBetween EachList failed:", err)
		return nil, err
	}

	log.Printf("IdentifyOwnerBetween %d issues found in list:\n", len(issues))

	retList := list.New()
	for _, issue := range issues {
		owner, err := identifyOwnerOne(issue, shortOnwers)

		if err == nil && owner.LastModify != "" {
			fmtModify, _ := utils.FormatTime2Day(owner.LastModifyTime)

			if !fmtModify.Before(startDay) && !fmtModify.After(endDay) {
				retList.PushBack(owner)
			} else {
				log.Printf("IdentifyOwnerBetween %d issues out of date:+%s\n", issue.Id, owner.LastModifyTime)
			}
		}
	}

	return retList, nil
}

func IdentifyOwnerOne(id int64, shortOnwers *list.List) (*IssueOwner, error) {
	issue, err := mantis.List(id)
	if err != nil {
		log.Print("IdentifyOwnerOne List failed:", err)
		return nil, err
	}

	return identifyOwnerOne(issue, shortOnwers)
}

func identifyOwnerOne(issue *mantis.Issue, shortOnwers *list.List) (*IssueOwner, error) {
	var err error

	var detail *mantis.IssueDetail
	detail, err = mantis.GetDetail(issue.Id)
	if err != nil {
		err = fmt.Errorf("Issue [%d] has no details", issue.Id)
		log.Print(err)
		return nil, err
	} else {
		owner := new(IssueOwner)
		owner.Id = issue.Id
		owner.Project = detail.Project
		owner.Level = issue.Level
		owner.Status = issue.Status
		owner.Summary = issue.Summary

		if DEBUG {
			log.Printf("Issue History [%d]:", owner.Id)
		}

		var history []mantis.IssueHistory
		err = json.Unmarshal([]byte(detail.History), &history)
		if err != nil {
			err = fmt.Errorf("Issue [%d] History error: %s", issue.Id, detail.History)
			log.Print(err)
			return nil, err
		}

		changeReg := regexp.MustCompile(`=>(.*)`)

		owneIssue := false

		for _, modify := range history {
			if DEBUG {
				log.Printf("	%s:%s	%s\n", modify.Username, modify.Field, modify.Change)
			}

			//LastModify
			if isInOwners(shortOnwers, modify.Username) {
				owner.LastModify = modify.Username
				owner.LastModifyTime = modify.DateModified
			}

			//LastAssignOutTo
			if modify.Field == "AssignedTo" {
				//log.Printf("Assigned To:%s => %q",modify.Change,assignedOugReg.FindStringSubmatch(modify.Change))

				match := changeReg.FindStringSubmatch(modify.Change)
				if len(match) < 2 {
					log.Printf("Issue [%d] History error: %s", issue.Id, detail.History)
				} else {
					newOwner := match[1]

					if isInOwners(shortOnwers, newOwner) { //分入
						owneIssue = true
						owner.LastAssignOutTo = ""
						if owner.InTime.IsZero() {
							owner.InTime = modify.DateModified
						}
					} else { //分出
						if owneIssue { //只有在分入状态才分出
							owner.LastAssignOutTo = newOwner
							owner.OutTime = modify.DateModified
							owneIssue = false
						}
					}

					if DEBUG {
						log.Printf("	newOwner=%s LastAssignOut=%s\n", newOwner, owner.LastAssignOutTo)
					}
				}
			}

			if modify.Field == "Status" {
				//log.Printf("Status To:%s => %q",modify.Change, changeReg.FindStringSubmatch(modify.Change))

				match := changeReg.FindStringSubmatch(modify.Change)
				if len(match) < 2 {
					log.Printf("Error Status Changes:%s => %q", modify.Change, changeReg.FindStringSubmatch(modify.Change))
				} else {
					if match[1] == "Readytorelease" || match[1] == "SCCBReview" {
						owner.LastFix = modify.Username

						//是组内人员解决的Bug
						if isInOwners(shortOnwers, modify.Username) {
							owner.OutTime = modify.DateModified
						}
					}
				}
			}
		}

		return owner, nil
	}
}

func isInOwners(shortOnwers *list.List, longName string) bool {
	//shortOnwers里存的是简单名词，如XXX，而Username可能是全称(xxx部门_xxx_139xxxx)，两者的关系不明
	for iter := shortOnwers.Front(); iter != nil; iter = iter.Next() {
		username := iter.Value.(string)

		if strings.Contains(longName, username) {
			return true
		}
	}

	return false
}
