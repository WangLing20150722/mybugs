package mantis

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"log"
	"sync"
)

var DEBUG = true

var instance *http.Client
var once sync.Once

func HTTPInstance() (*http.Client) {
	once.Do(func() {
		cookieJar, _ := cookiejar.New(nil)
		instance = &http.Client{Jar: cookieJar}
	})

	return instance
}

func DumpCookies() {
	client := HTTPInstance()

	curl,err := url.Parse("http://mantis.tclking.com/")
	if(err != nil) {
		log.Fatalf("url.Parse failed %s\n",err)
		return
	}

	cookies := client.Jar.Cookies(curl)
	log.Print(cookies)
	for _, cookie := range cookies {
		log.Print(cookie)
	}
}

