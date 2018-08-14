package mantis

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"time"
)

type IssueDetail struct {
	Id        int64     `gorm:"primary_key;auto_increment:false"`
	FetchTime time.Time `gorm:"column:page_created"`
	History   string    `gorm:"type:text"`
	Project   string    `gorm:"type:text"`
}

type IssueHistory struct {
	DateModified time.Time
	Username     string
	Field        string
	Change       string //状态改变
}

func GetIssueDetail(id int64) (*IssueDetail, error) {
	client := HTTPInstance()

	surl := fmt.Sprintf(`http://mantis.tclking.com/view.php?id=%d`, id)

	if DEBUG {
		log.Printf("GetIssueDetail id=%d url=%s\n", id, surl)
	}

	req, err := http.NewRequest("GET", surl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	now := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("GetIssueDetail id=%d StatusCode = %d\n", id, resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetIssueDetail id=%d StatusCode = %d\n", id, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	if DEBUG {
		log.Printf("GetIssueDetail id %d Response:\n", id)
		log.Println(doc.Html())
	}

	sel := doc.Find("#history_open > table > tbody > tr")
	if sel.Length() <= 0 {
		return nil, fmt.Errorf("GetIssueDetail error page,can not find history_open")
	}

	issueDetail := new(IssueDetail)
	issueDetail.FetchTime = now
	issueDetail.Id = id

	var historys []*IssueHistory

	sel.Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		if tds.Length() <= 0 {
			log.Printf("Find td %d error", i)
			return
		}

		history := new(IssueHistory)

		tds.Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0: //Date Modified
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error Date Modified 1 %d %s", i, html)
				}
				value = strings.Replace(value, " ", "", -1)
				value = strings.Replace(value, "\n", "", -1)
				value = strings.Replace(value, "\t", "", -1)

				var err error
				history.DateModified, err = time.Parse("2006-01-0215:04", value)
				if err != nil {
					log.Printf("Error Date Modified 1 %d [%s] %s", i, value, err)
				}
				break

			case 1: //Username
				value := s.Children().Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error Username %d %s", i, html)
				}
				value = strings.Replace(value, " ", "", -1)
				value = strings.Replace(value, "\n", "", -1)
				value = strings.Replace(value, "\t", "", -1)
				history.Username = value
				break

			case 2: //Field
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error Field %d %s", i, html)
				}
				value = strings.Replace(value, " ", "", -1)
				value = strings.Replace(value, "\n", "", -1)
				value = strings.Replace(value, "\t", "", -1)
				history.Field = value
				break

			case 3: //Change
				value := s.Text()
				if value == "" {
					html, _ := s.Html()
					log.Printf("Error Change %d %s", i, html)
				}
				value = strings.Replace(value, " ", "", -1)
				value = strings.Replace(value, "\n", "", -1)
				value = strings.Replace(value, "\t", "", -1)
				history.Change = value
				break
			}
		})

		historys = append(historys, history)
	})

	b, err := json.Marshal(historys)
	if err == nil {
		issueDetail.History = string(b[:])
	} else {
		log.Printf("GetIssueDetail Marshal historys error:%v\n", err)
	}

	projectSel := doc.Find("#view-issue-details > table > tbody > tr.bug-header-data > td.bug-project")
	if projectSel.Length() > 0 {
		issueDetail.Project = projectSel.Text()
	} else {
		log.Printf("GetIssueDetail Marshal historys error:%v\n", err)
	}

	return issueDetail, nil
}
