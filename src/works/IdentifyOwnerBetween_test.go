package works_test

import (
	"container/list"
	"fmt"
	"mantis"
	"testing"
	"time"
	"works"
)

func init() {
	mantis.OpenDB()
}

func TestIdentifyOwnerBetweenOneDay(t *testing.T) {
	startday, _ := time.Parse("2006-01-02", "2018-04-01")
	endday, _ := time.Parse("2006-01-02", "2018-04-02")

	l, err := works.IdentifyOwnerBetween(startday, endday, getShortNameList())
	if err != nil {
		t.Fatal(err)
	}

	works.IssueOwnerList2Excel(l, "test.xlsx")

	for iter := l.Front(); iter != nil; iter = iter.Next() {
		issue := iter.Value.(*works.IssueOwner)
		fmt.Printf("doTestIdentifyOwnerOne(t,%d,\"%s\",\"%s\",\"%s\")\n", issue.Id, issue.LastModify, issue.LastAssignOutTo, issue.LastFix)
	}
}

func TestIdentifyOwnerOnes(t *testing.T) {
	doTestIdentifyOwnerOne(t, 174082, "电视应用所_曾辉云_13510630237", "中国软件_傅钇中_15986521123", "中国软件_傅钇中_15986521123")
	doTestIdentifyOwnerOne(t, 170157, "电视应用所_曹绍华_15989331339", "中国软件_曾泽欢_15118100132", "中国软件_曾泽欢_15118100132")
	doTestIdentifyOwnerOne(t, 164679, "电视应用所_曹绍华_15989331339", "", "电视应用所_曹绍华_15989331339")
	doTestIdentifyOwnerOne(t, 170266, "电视应用所_曹绍华_15989331339", "", "电视应用所_曹绍华_15989331339")
	doTestIdentifyOwnerOne(t, 166842, "电视应用所_曹绍华_15989331339", "中国软件所_杨健_18428399071", "中国软件_曾泽欢_15118100132")
	doTestIdentifyOwnerOne(t, 173531, "电视应用所_包林玉_15889484305", "研发中心西安软件所_霍平华_18629302436", "")
	doTestIdentifyOwnerOne(t, 174374, "电视应用所_周孙强_18932470732", "", "电视应用所_周孙强_18932470732")
	doTestIdentifyOwnerOne(t, 165043, "电视应用所_曾辉云_13510630237", "中国软件_傅钇中_15986521123", "TV+终端所_唐鹏_17097218985")
}

func getShortNameList() *list.List {
	l := list.New()
	l.PushBack("包林玉")
	l.PushBack("曾辉云")
	l.PushBack("曹绍华")
	l.PushBack("周孙强")
	return l
}

func doTestIdentifyOwnerOne(t *testing.T, id int64, LastModify, LastAssignOutTo, LastFix string) {
	owner, err := works.IdentifyOwnerOne(id, getShortNameList())
	if err != nil {
		t.Fatal(err)
	}
	if owner.Id != id {
		t.Fatal("id error ", id, owner.Id)
	}
	if owner.LastModify != LastModify {
		t.Fatal("LastModify error", owner.LastModify, LastModify)
	}
	if owner.LastAssignOutTo != LastAssignOutTo {
		t.Fatal("LastAssignOutTo error", owner.LastAssignOutTo, LastAssignOutTo)
	}

	if owner.LastFix != LastFix {
		t.Fatal("LastFix error", owner.LastFix, LastFix)
	}
}
