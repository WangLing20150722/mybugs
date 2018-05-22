package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type CONFIG_S struct {
	Username string
	Password string
}

var CONFIG CONFIG_S

func init() {
	//For testcase, save username & password in password.json
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		filename := os.Getenv("GOPATH") + "/password.json"

		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("ReadNameList ReadFile error:%s", err)
			return
		}

		err = json.Unmarshal(bytes, &CONFIG)
		if err != nil {
			log.Printf("ReadNameList Unmarshal error:%s", err)
			return
		}

		log.Print("Read password from password.json for test:", CONFIG)
	}
}
