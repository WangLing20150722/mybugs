package main

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	err = c.Run(ctxt, googleSearch("site:brank.as", "Home", &site, &res))
	if err != nil {
		log.Fatal(err)
	}

	/*
	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
	*/

	log.Printf("saved screenshot from search result listing `%s` (%s)", res, site)
}

func googleSearch(q, text string, site, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		//open home page
		chromedp.Navigate(`http://mantis.tclking.com/`),

		//login
		chromedp.WaitVisible(`#login-form`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>> Open homepage done!")
			return nil
		}),
		chromedp.SetValue(`#username`, "lihui02" ,chromedp.ByID),
		chromedp.SetValue(`#password`, "asdfzxcv" ,chromedp.ByID),
		chromedp.SetValue(`#remember-login`, "asdfzxcv" ,chromedp.ByID),
		chromedp.SetAttributeValue(`#remember-login`,"checked","true",chromedp.ByID),
		chromedp.RemoveAttribute(`#secure-session`,"checked",chromedp.ByID),
		chromedp.Click(`.button`,chromedp.ByQuery),

		//select project
		chromedp.WaitVisible(`#form-set-project-id`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>> Login done")
			return nil
		}),
		chromedp.SetJavascriptAttribute(`#form-set-project-id`,"selectedIndex","0",chromedp.ByID),
		chromedp.Submit(`#form-set-project`,chromedp.ByQuery),
		//chromedp.Click(`#form-set-project`,chromedp.ByQuery),

		//view all bugs
		chromedp.WaitVisible(`#menu-items`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>> Project selected done")
			return nil
		}),
		chromedp.Navigate(`http://mantis.tclking.com/view_all_bug_page.php`),

		//filter bugs
		chromedp.WaitVisible(`#menu-items`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>> view all bugs done")
			return nil
		}),





	}
}