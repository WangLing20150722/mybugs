package mantis

import (
	"container/list"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var s_dbfile = "issues.db"

func ClearList() error {
	return os.Remove(s_dbfile)
}

func SaveList(l *list.List) error {
	db, err := gorm.Open("sqlite3", s_dbfile)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&Issue{})

	for e := l.Front(); e != nil; e = e.Next() {
		if(db.NewRecord(e.Value.(*Issue))){
			db.Create(e.Value.(*Issue))
		} else {
			db.Save(e.Value.(*Issue))
		}
	}

	return nil
}

func SaveDetail(detail *IssueDetail) error {
	db, err := gorm.Open("sqlite3", s_dbfile)
	if err != nil {
		return err
	}
	defer db.Close()

	db.AutoMigrate(&IssueDetail{})

	if(db.NewRecord(detail)) {
		db.Create(detail)
	} else {
		db.Save(detail)
	}

	return nil
}
