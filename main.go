package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/huandu/facebook"
)

var (
	appID = "<app_id>"
	appSecret = "<app_secret>"
	profileID = "310028155725467" // Electoral Commission Ghana Page
	pattern = "Presidential Provisional Results"
	session *facebook.Session
)

func main() {
	fb := facebook.New(appID, appSecret)
	token := fb.AppAccessToken()
	log.Print("Got token: ", token)

	session = fb.Session(token)

	readFeed(profileID, token)
}

func readFeed(profileID, token string) {
	url := fmt.Sprintf("/%s/feed", profileID)
	res, err := session.Get(url, facebook.Params{
		"access_token": token,
		"since": "1481155200",
		"limit": "100",
	})

	if err != nil {
		log.Fatal(err)
	}

	var items []facebook.Result
	err = res.DecodeField("data", &items)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, item := range items {
		msg := fmt.Sprintf("%s",item["message"])
		match, _ := regexp.MatchString(pattern, msg)
		if match {
			getData(msg)
		}
	}
}

func getData(msg string) {
	distre := regexp.MustCompile("/*.District: (<?P<dist>.*)/i")
	
	if distre.MatchString(msg) {
		log.Print("District: ", distre.SubexpNames()[1])
	} else {
		log.Print("does not match")
	}
}

func searchPage(token string) {
	res, _ := facebook.Get("/search", facebook.Params{
		"access_token": token,
		"type":         "page",
		"q":            "Electoral Commission Ghana",
	})

	var items []facebook.Result
	err := res.DecodeField("data", &items)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fmt.Println(item)
	}
}
