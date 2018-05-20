package mantis

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

var s_projectid = ""
var s_projectName = ""

func ProjectId() string {
	return s_projectid
}

func ProjectName() string {
	return s_projectName
}

func SetProjectById(id string) (bool, error) {
	s_projectid = ""
	s_projectName = ""

	if !CheckLogin(nil) {
		return false, fmt.Errorf("SetProjectById not login")
	}

	client := HTTPInstance()

	post := fmt.Sprintf("project_id=%s", id)
	if DEBUG {
		log.Printf("SetProjectById post=%s\n", post)
	}

	form := strings.NewReader(post)
	req, err := http.NewRequest("POST", "http://mantis.tclking.com/set_project.php", form)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	log.Printf("SetProjectById StatusCode = %d\n", resp.StatusCode)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("SetProjectById failed StatusCode=%d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false, err
	}

	if DEBUG {
		log.Println("SetProjectById Response:")
		log.Println(doc.Html())
	}

	if !CheckLogin(doc) {
		return false, fmt.Errorf("SetProjectById not login after request")
	}

	ret := VerifyProjectId(doc, id)
	log.Printf("SetProjectById to %s-%s\n", s_projectid, s_projectName)

	if ret {
		return true, nil
	} else {
		return false, fmt.Errorf("SetProjectById failed want %s but got %s", id, s_projectid)
	}
}

func VerifyProjectId(doc *goquery.Document, projectid string) bool {
	if projectid == "" {
		projectid = s_projectid
	}

	sel := doc.Find("#form-set-project-id")
	if len(sel.Nodes) <= 0 {
		return false
	}

	ret := false
	sel.Children().EachWithBreak(func(i int, s *goquery.Selection) bool {
		value, _ := s.Attr("value")
		_, isselected := s.Attr("selected")

		if isselected {
			if value == projectid {
				ret = true
			} else {
				ret = false
				log.Fatalf("Error: Project wants to be %s but be %s\n", projectid, value)
			}

			s_projectid = value
			s_projectName = s.Text()

			return false
		}

		return true
	})

	return ret
}
