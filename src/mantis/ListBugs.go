package mantis

import (
	"container/list"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Issue struct {
	Id            int64     `gorm:"primary_key;auto_increment:false"`
	Level         string    `gorm:"type:text"`
	Summary       string    `gorm:"type:text"`
	Status        string    `gorm:"type:text"`
	Reporter      string    `gorm:"type:text"`
	AssignedTo    string    `gorm:"type:text"`
	DateSubmitted string    `gorm:"type:text"`
	DueDate       string    `gorm:"type:text"`
	Updated       string    `gorm:"type:text"`
	FetchTime     time.Time `gorm:"column:page_created"`
}

var s_sort = ""
var s_asc = false

var s_statusMap = map[string]string{
	"S": "SCCB Review",
	"C": "Closed",
	"A": "Assigned",
	"R": "Ready to release",
	"M": "More info requested",
	"N": "New",
	"F": "Failed",
	"P": "Pretest",
	"T": "Ready to test",
	"D": "Ready to deploy",
}

/**
List All Bugs by page,sort and order.
sort - last_updated
*/
func ListBugs(page int, sort string, asc bool) (*list.List, error) {
	l := list.New()

	now := time.Now()

	doc, err := sortIssues(sort, asc)
	if err != nil {
		return l, err
	}

	if page != 0 {
		doc, err = openPage(page, sort, asc)
	}
	if err != nil {
		return l, err
	}

	if !parseBugList(doc, l, now) {
		return l, fmt.Errorf("ListBugs parseBugList %d failed", page)
	}

	return l, nil
}

func sortIssues(sort string, asc bool) (*goquery.Document, error) {
	if s_sort == sort && asc == s_asc {
		//log.Printf("SortIssues already sort by %s and asc=%t",sort,asc)
		return nil, nil
	}

	client := HTTPInstance()

	var dir string
	if asc {
		dir = "ASC"
	} else {
		dir = "DESC"
	}

	surl := fmt.Sprintf("http://mantis.tclking.com/view_all_set.php?sort=%s&dir=%s&type=2", sort, dir)

	if DEBUG {
		log.Printf("SortIssues url=%s\n", surl)
	}
	req, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("SortIssues StatusCode = %d\n", resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ListBugs failed StatusCode=%d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if DEBUG {
		log.Println("Login Response:")
		log.Println(doc.Html())
	}

	if !CheckLogin(doc) {
		return nil, fmt.Errorf("ListBugs Login failed")
	}

	if !VerifyProjectId(doc, "") {
		return nil, fmt.Errorf("ListBugs VerifyProjectId failed")
	}

	if !checkOrder(doc, sort, asc) {
		return nil, fmt.Errorf("ListBugs set order failed")
	}

	s_sort = sort
	s_asc = asc

	return doc, nil
}

func checkOrder(doc *goquery.Document, sort string, asc bool) bool {
	var query string
	if asc {
		query = `img[src="http://mantis.tclking.com/images/up.gif"]`
	} else {
		query = `img[src="http://mantis.tclking.com/images/down.gif"]`
	}

	sel := doc.Find(query)
	if len(sel.Nodes) <= 0 {
		return false
	}

	ret := false
	strsort := "sort=" + sort

	sel.Each(func(i int, s *goquery.Selection) {
		href, has := s.Prev().Attr("href")

		if has && strings.Contains(href, strsort) {
			ret = true
		}
	})

	return ret
}

func openPage(page int, sort string, asc bool) (*goquery.Document, error) {
	client := HTTPInstance()

	surl := fmt.Sprintf(`http://mantis.tclking.com/view_all_bug_page.php?page_number=%d`, page+1)

	if DEBUG {
		log.Printf("ListBugs openPage=%d url=%s\n", page, surl)
	}

	req, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("ListBugs openPage %d StatusCode = %d\n", page, resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ListBugs openPage %d failed StatusCode=%d", page, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if DEBUG {
		log.Printf("ListBugs openPage %d Response:\n", page)
		log.Println(doc.Html())
	}

	if !CheckLogin(doc) {
		return nil, fmt.Errorf("ListBugs openPage %d Login failed", page)
	}

	if !VerifyProjectId(doc, "") {
		return nil, fmt.Errorf("ListBugs openPage %d VerifyProjectId failed", page)
	}

	if !checkOrder(doc, sort, asc) {
		return nil, fmt.Errorf("ListBugs openPage %d set order failed", page)
	}

	return doc, nil
}

func parseBugList(doc *goquery.Document, l *list.List, now time.Time) bool {
	sel := doc.Find("#buglist > tbody > tr")
	if sel.Length() <= 0 {
		log.Printf("parseBugList error page")
		return false
	}

	sel.Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		if class == "buglist-footer" {
			return
		}

		tds := s.Find("td")
		if tds.Length() <= 0 {
			log.Printf("Find td %d error", i)
			return
		}

		issue := new(Issue)
		issue.FetchTime = now

		tds.Each(func(i int, s *goquery.Selection) {
			class, has := s.Attr("class")
			if !has {
				html, _ := s.Html()
				log.Printf("Error issue no class %d %s", i, html)
			}
			switch class {
			case "column-id":
				value := s.Children().Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				var err error
				issue.Id, err = strconv.ParseInt(value, 10, 0)
				if err != nil {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				break

			case "column-severity":
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				issue.Level = value
				break

			case "column-summary":
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				issue.Summary = value
				break

			case "column-status":
				value := s.Children().Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				status, has := s_statusMap[value]
				if has {
					issue.Status = status
				} else {
					log.Printf("Error issue %d unknown status %s", i, value)
					issue.Status = value
				}

				break

			case "column-reporter":
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				issue.Reporter = value
				break

			case "column-assigned-to":
				value := s.Children().Text()
				issue.AssignedTo = value
				break

			case "column-date-submitted":
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				issue.DateSubmitted = value
				break

			case "column-due-date":
				value := s.Text()
				issue.DueDate = value
				break

			case "column-last-modified":
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error issue %d %s", i, html)
				}
				issue.Updated = value
				break
			}
		})

		l.PushBack(issue)
	})

	return true
}
