package mantis

import (
	"strings"
	"net/http"
	"fmt"
	"log"
	"github.com/PuerkitoBio/goquery"
)

var login = false

func IsLogin() bool {
	return login
}

func Login(username,password string) (bool,error) {
	login = false

	client := HTTPInstance()

	post := fmt.Sprintf("return=view_all_bug_page.php&username=%s&password=%s&perm_login=on",username,password)
	if(DEBUG) {
		log.Printf("Login post=%s\n",post)
	}

	form := strings.NewReader(post)
	req, err := http.NewRequest("POST","http://mantis.tclking.com/login.php",form)
	if err != nil {
		return false,err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	log.Printf("Login StatusCode = %d\n",resp.StatusCode)
	defer resp.Body.Close()

	if(resp.StatusCode != http.StatusOK) {
		return false,fmt.Errorf("Login failed StatusCode=%d", resp.StatusCode)
	}


	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false,err
	}

	if(DEBUG) {
		log.Println("Login Response:")
		log.Println(doc.Html())
	}

	if( _checkLogin(doc) ) {
		login = true
		log.Println("Login Sucess")
		return true,nil
	} else {
		return false,fmt.Errorf("Login failed")
	}
}

func CheckLogin(doc* goquery.Document) bool {
	if(!login) {
		return false
	}

	if doc == nil {
		return true
	}

	return _checkLogin(doc)
}

func _checkLogin(doc* goquery.Document) bool {
	if( len(doc.Find("#login-info").Nodes) > 0 ) {
		return true
	} else {
		return false
	}
}