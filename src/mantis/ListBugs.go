package mantis

import (
	"fmt"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
	"container/list"
	"strconv"
)

type Issue struct {
	Id 		int64	`gorm:"primary_key;auto_increment:false"`
	Level   	string	`gorm:"type:text"`
	Summary 	string	`gorm:"type:text"`
	Status  	string	`gorm:"type:text"`
	Reporter 	string	`gorm:"type:text"`
	AssignedTo 	string	`gorm:"type:text"`
	DateSubmitted	string	`gorm:"type:text"`
	DueDate		string	`gorm:"type:text"`
	Updated		string	`gorm:"type:text"`
}

/**
List All Bugs by page,sort and order.
sort - last_updated
 */
func ListBugs(pageoffset,pagelimited int, sort string, asc bool) (*list.List,error) {
	client := HTTPInstance()

	time.Now()

	var dir string
	if(asc) {
		dir = "ASC"
	} else {
		dir = "DESC"
	}

	surl := fmt.Sprintf("http://mantis.tclking.com/view_all_set.php?sort=%s&dir=%s&type=2",sort,dir)

	if(DEBUG) {
		log.Printf("ListBugs pageoffset=%d pagelimited=%d url=%s\n",pageoffset,pagelimited,surl)
	}
	req, err := http.NewRequest("GET",surl,nil)
	if err != nil {
		return nil,err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	log.Printf("ListBugs StatusCode = %d\n",resp.StatusCode)
	defer resp.Body.Close()

	if(resp.StatusCode != http.StatusOK) {
		return nil,fmt.Errorf("ListBugs failed StatusCode=%d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil,err
	}

	if(DEBUG) {
		log.Println("Login Response:")
		log.Println(doc.Html())
	}

	if( !CheckLogin(doc) ) {
		return nil,fmt.Errorf("ListBugs Login failed")
	}

	if(!VerifyProjectId(doc,"")) {
		return  nil,fmt.Errorf("ListBugs VerifyProjectId failed")
	}

	if(!checkOrder(doc,sort,asc)) {
		return nil,fmt.Errorf("ListBugs set order failed")
	}

	l := list.New()
	for i:=pageoffset;i<pageoffset+pagelimited;i++ {
		if(i!=0) {
			doc,err = openPage(i,sort,asc)
			if(err != nil) {
				return l,err
			}
		}

		if(!parseBugList(doc,l)) {
			return l,fmt.Errorf("ListBugs parseBugList %d failed",i)
		}
	}

	return l,nil
}

func checkOrder(doc* goquery.Document, sort string, asc bool) bool {
	var query string
	if(asc) {
		query = `img[src="http://mantis.tclking.com/images/up.gif"]`
	} else {
		query = `img[src="http://mantis.tclking.com/images/down.gif"]`
	}

	sel := doc.Find(query)
	if(len(sel.Nodes) <= 0) {
		return false;
	}

	ret := false
	strsort := "sort=" + sort

	sel.Each(func (i int,s *goquery.Selection ) {
		href, has := s.Prev().Attr("href")

		if (has && strings.Contains(href, strsort)) {
			ret = true
		}
	})

	return ret
}

func openPage(page int,sort string, asc bool) (*goquery.Document,error) {
	client := HTTPInstance()

	surl := fmt.Sprintf(`http://mantis.tclking.com/view_all_bug_page.php?page_number=%d`,page+1)

	if(DEBUG) {
		log.Printf("ListBugs openPage=%d url=%s\n",page,surl)
	}

	req, err := http.NewRequest("GET",surl,nil)
	if err != nil {
		return nil,err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	log.Printf("ListBugs openPage %d StatusCode = %d\n",page,resp.StatusCode)
	defer resp.Body.Close()

	if(resp.StatusCode != http.StatusOK) {
		return nil,fmt.Errorf("ListBugs openPage %d failed StatusCode=%d", page,resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil,err
	}

	if(DEBUG) {
		log.Printf("ListBugs openPage %d Response:\n",page)
		log.Println(doc.Html())
	}

	if( !CheckLogin(doc) ) {
		return nil,fmt.Errorf("ListBugs openPage %d Login failed",page)
	}

	if(!VerifyProjectId(doc,"")) {
		return  nil,fmt.Errorf("ListBugs openPage %d VerifyProjectId failed",page)
	}

	if(!checkOrder(doc,sort,asc)) {
		return nil,fmt.Errorf("ListBugs openPage %d set order failed",page)
	}

	return doc,nil
}

func parseBugList(doc* goquery.Document,l *list.List) bool {
	sel := doc.Find("#buglist > tbody > tr")
	if(sel.Length() <= 0) {
		log.Printf("parseBugList error page")
		return false;
	}

	sel.Each(func (i int,s *goquery.Selection ) {
		class,_ := s.Attr("class")
		if(class == "buglist-footer") {
			return
		}

		tds := s.Find("td")
		if(tds.Length() <= 0) {
			log.Printf("Find td %d error",i)
			return
		}

		issue := new(Issue)

		tds.Each(func (i int,s *goquery.Selection ) {
			class,has := s.Attr("class")
			if(!has ) {
				html,_ := s.Html()
				log.Printf("Error issue no class %d %s",i,html)
			}
			switch(class) {
			case "column-id":
				value := s.Children().Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				var err error
				issue.Id,err = strconv.ParseInt(value,10,0)
				if(err != nil) {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				break

			case "column-severity":
				value := s.Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.Level = value
				break

			case "column-summary":
				value := s.Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.Summary = value
				break

			case "column-status":
				value,exist := s.Children().Attr("title")
				if(!exist) {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.Status = value
				break

			case "column-reporter":
				value := s.Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.Reporter = value
				break

			case "column-assigned-to":
				value := s.Children().Text()
				issue.AssignedTo = value
				break

			case "column-date-submitted":
				value := s.Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.DateSubmitted = value
				break

			case "column-due-date":
				value := s.Text()
				issue.DueDate = value
				break

			case "column-last-modified":
				value := s.Text()
				if(value == "") {
					html,_ := s.Html()
					log.Printf("Error issue %d %s",i,html)
				}
				issue.Updated = value
				break
			}
		});

		l.PushBack(issue)
	})

	return true
}
