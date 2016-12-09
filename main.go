package main

import (
	"fmt"
	"log"
	"regexp"
	"os"
	"text/template"
	"strings"
	"strconv"

	"github.com/huandu/facebook"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

type entry struct {
	ID string `bson:"_id"`
	District string
	Constituency string
	CPP int
	NDP int
	NDC int
	PPP int
	NPP int
	PNC int
	IND int
	Total int
	Rejected int
	InBox int
}

var (
	appID = "<app_id>"
	appSecret = "<app_secret>"
	profileID = "310028155725467" // Electoral Commission Ghana Page
	pattern = "Presidential Provisional Results"
	mongohost = "localhost"

	session *facebook.Session
	collection *mgo.Collection
)

func main() {
	log.Print("Checking for new posts ...")

	fb := facebook.New(appID, appSecret)
	token := fb.AppAccessToken()
	session = fb.Session(token)

	dbsession, err := mgo.Dial(mongohost)
        if err != nil {
                panic(err)
        }
        defer dbsession.Close()

        // Optional. Switch the session to a monotonic behavior.
        dbsession.SetMode(mgo.Monotonic, true)
	collection = dbsession.DB("ecresults").C("results")

	n := readFeed(profileID, token)

	log.Print("Got ", n, " declarations.")
}

func readFeed(profileID, token string) int {
	count := 0

	csvheader := "District,Constituency,CPP,NDP,NDC,PPP,NPP,PNC,IND,Total,Rejected,In-Box"

	fmt.Println(csvheader)

	url := fmt.Sprintf("/%s/feed", profileID)
	res, err := session.Get(url, facebook.Params{
		"access_token": token,
		"since": "1481155200", 	// 2016-12-08
		"until": "1481846400",	// 2016-12-16
		"limit": "100",
	})

	if err != nil {
		log.Fatal(err)
	}

	paging, _ := res.Paging(session)
	items := paging.Data()

	row := &entry{}

	for ;; {
		items = paging.Data()
		for _, item := range items {
			msg := fmt.Sprintf("%s",item["message"])
			row.ID = fmt.Sprintf("%s", item["id"])
			match, _ := regexp.MatchString(pattern, msg)
			if match {
				getData(msg, row)
				count++
			}
		}

		if paging.HasNext() {
			paging.Next()
		} else {
			break
		}
	}

	return count
}

func getData(msg string, row *entry) {

	const line = `{{.District}},{{.Constituency}},{{.CPP}},{{.NDP}},{{.NDC}},{{.PPP}},{{.NPP}},{{.PNC}},{{.IND}},{{.Total}},{{.Rejected}},{{.InBox}}
`
	tmpl := template.Must(template.New("line").Parse(line))

	re := regexp.MustCompile(`Name of District: (?P<dist>[\w ]*)`)
	match := re.FindStringSubmatch(msg)
	row.District = match[1]

	re = regexp.MustCompile(`Name of Constituency: (?P<const>[\w ]*)`)
	match = re.FindStringSubmatch(msg)
	row.Constituency = match[1]

	re = regexp.MustCompile(`CPP:\s*(?P<cpp>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.CPP, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`NDP:\s*(?P<ndp>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.NDP, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`NDC:\s*(?P<ndc>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.NDC, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))
	
	re = regexp.MustCompile(`PPP:\s*(?P<ppp>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.PPP, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`NPP:\s*(?P<npp>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.NPP, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`PNC:\s*(?P<pnc>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.PNC, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`IND:\s*(?P<ind>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.IND, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`Total Votes:\s*(?P<total>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.Total, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`Rejected Votes:\s*(?P<rejected>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.Rejected, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	re = regexp.MustCompile(`Total Votes in Ballot Box:\s*(?P<inbox>[\d,]*)`)
	match = re.FindStringSubmatch(msg)
	row.InBox, _ = strconv.Atoi(strings.Replace(match[1], ",", "", -1))

	_, err := collection.UpsertId(row.ID, row)
	if err != nil {
		log.Fatal(err)
	}

	// output as CSV to stdout
	err = tmpl.Execute(os.Stdout, row)
	if err != nil {
		log.Fatal(err)
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
