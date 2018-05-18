package mantis

import (
	"container/list"
	"os"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func SaveList(l *list.List,dbfile string) error {
	os.Remove(dbfile)

	db, err := gorm.Open("sqlite3", dbfile)
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
