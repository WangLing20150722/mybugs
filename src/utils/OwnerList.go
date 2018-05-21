package utils

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func GenerateNameList() {
	owners := [...]string{
		"电视应用所_包林玉_15889484305",
		"电视应用所_欧阳英俊_18688967602",
		"电视应用所_卢美辉",
		"电视应用所_张俊_18823891034",
		"电视应用所_周孙强_18932470732",
		"电视应用所_刘嘉宁_15217769295",
		"电视应用所_杨建军_18589027903",

		"电视应用所_杨少东_075533910215",

		"电视应用所_曾辉云_13510630237",
		"电视应用所_吴晓红_13684922100",
		"电视应用所_黄孙明_13692224232",
		"电视应用所_郝子聪_13823351326",
		"电视应用所_陈闯",
		"电视应用所_王明月_15619285760",
		"电视应用所_王芳_13530992896",
		"杨倩倩",

		"电视应用所_刘远凤_15817468435",
		"电视应用所_曹绍华_15989331339",

		"电视应用所_刘礼丰_17098870199",
		"电视应用所_李章平_13537541287",
		"电视应用所_李蒙蒙_15914197631",
		"电视应用所_欧阳雪琴_18810732928",
		"电视应用所_洪东金_18740441230",
		"商鑫奇",

		"电视应用所_罗克光_18824290329",
		"电视应用所_熊伟",
		"电视应用所_袁胜_18291906071",
		"电视应用所_赖国锋_13923761987",
		"刘二虎",
		"万旭",
		"李志纯",
		"蔡祥林",
		"黄振宇",
		"李棚",
		"田园",
		"魏炯",
	}

	js, _ := json.MarshalIndent(owners, "", "")
	fmt.Printf(string(js))
}

func ReadNameList(filename string) *list.List {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ReadNameList ReadFile error:%s", err)
		return nil
	}

	var owners []string
	err = json.Unmarshal(bytes, &owners)
	if err != nil {
		log.Printf("ReadNameList Unmarshal error:%s", err)
		return nil
	}

	l := list.New()
	for _, owner := range owners {
		l.PushBack(owner)
	}

	return l
}
