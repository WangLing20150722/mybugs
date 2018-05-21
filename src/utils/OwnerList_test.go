package utils_test

import (
	"log"
	"testing"
	"utils"
)

func TestGenerateNameList(t *testing.T) {
	utils.GenerateNameList()
}

func TestReadNameList(t *testing.T) {
	l := utils.ReadNameList("../../owners.json")

	if l == nil {
		t.Fail()
	}

	if l.Len() <= 0 {
		t.Fatal()
	}

	for iter := l.Front(); iter != nil; iter = iter.Next() {
		log.Print(iter.Value.(string))
	}
}
