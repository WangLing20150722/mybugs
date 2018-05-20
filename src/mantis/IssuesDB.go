package mantis

import (
	"container/list"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"time"
)

var s_dbfile = "issues.db"
var s_db *gorm.DB

func OpenDB() error {
	if s_db != nil {
		return fmt.Errorf("DB already Openned")
	}
	var err error
	s_db, err = gorm.Open("sqlite3", s_dbfile)
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	if s_db != nil {
		return s_db.Close()
	}

	return fmt.Errorf("Db not Opened")
}

func ClearList() error {
	return os.Remove(s_dbfile)
}

func SaveList(l *list.List) error {
	if s_db == nil {
		return fmt.Errorf("Db not Opened")
	}

	s_db.AutoMigrate(&Issue{})

	for e := l.Front(); e != nil; e = e.Next() {
		if s_db.NewRecord(e.Value.(*Issue)) {
			s_db.Create(e.Value.(*Issue))
		} else {
			s_db.Save(e.Value.(*Issue))
		}
	}

	return nil
}

func EachList(start, end time.Time) ([]*Issue, error) {
	if s_db == nil {
		return nil, fmt.Errorf("Db not Opened")
	}

	var Issues []*Issue
	s_db.Find(&Issues, "Updated >= ? AND Updated <= ?", start.Format("2006-01-02"), end.Format("2006-01-02"))
	return Issues, nil
}

func List(id int64) (*Issue, error) {
	if s_db == nil {
		return nil, fmt.Errorf("Db not Opened")
	}

	var iss Issue
	s_db.First(&iss, "Id = ?", id)
	return &iss, nil
}

func SaveDetail(detail *IssueDetail) error {
	if s_db == nil {
		return fmt.Errorf("Db not Opened")
	}

	s_db.AutoMigrate(&IssueDetail{})

	if s_db.NewRecord(detail) {
		s_db.Create(detail)
	} else {
		s_db.Save(detail)
	}

	return nil
}

func GetDetail(id int64) (*IssueDetail, error) {
	if s_db == nil {
		return nil, fmt.Errorf("Db not Opened")
	}

	var detail IssueDetail
	s_db.First(&detail, "Id = ?", id)
	if detail.Id == id {
		return &detail, nil
	} else {
		return nil, fmt.Errorf("Issue(%d) not Found", id)
	}
}
