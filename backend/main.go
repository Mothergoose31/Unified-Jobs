package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

type RequestData struct {
	Title    string `json:"title"`
	Location string `json:"location"`
}

type Job struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Url      string `json:"url"`
}

var reqData = RequestData{Title: "", Location: ""}
var jobs = make([]Job, 0, 300)

//Job boardSites that do not allow webscrapping

//https://linkedin.com/robots.txt
//https://www.monster.com/robots.txt

//================================================================
// can be scrapped
//dice.com/jobs?q="

func scrapeIndeed(c *colly.Collector) {
	c.OnHTML("div.jobsearch-SerpJobCard", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.jobtitle"))
		company := strings.TrimSpace(e.ChildText("span.company"))
		location := strings.TrimSpace(e.ChildText("span.location"))
		url := strings.TrimSpace(e.ChildAttr("h2.title > a", "href"))

		if !strings.Contains(url, "https://indeed.com") {
			url = "https://indeed.com" + url
		}

		job := Job{
			Title:    title,
			Company:  company,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" && url != "" {
			jobs = append(jobs, job)
		}
	})
}

func scrapeStackOverFlow(c *colly.Collector) {
	c.OnHTML("div.-job", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText("a.s-link"))
		company := strings.TrimSpace(e.ChildText("h3 > span"))
		// Remove newline char and all chars after newline in company name
		companyFinal := company[0:strings.Index(company, "\n")]
		location := strings.TrimSpace(e.ChildText("span.fc-black-500"))
		url := strings.TrimSpace(e.ChildAttr("a.s-link[href]", "href"))
		url = "https://stackoverflow.com" + url

		job := Job{
			Title:    title,
			Company:  companyFinal,
			Location: location,
			Url:      url,
		}

		if location == "" {
			location = reqData.Location
		}
		if title != "" && url != "" {
			jobs = append(jobs, job)
		}
	})
}

func Handle(w http.ResponseWriter, r *http.Request) {
	//@dev get the request for the query along with
	//the parameters the title being searched and the location

	query := r.URL.Query()
	reqData.Title = query.Get("title")
	reqData.Location = query.Get("location")

	//  normalize search query terms for scraping
	strings.Replace(reqData.Title, " ", "+", -1)
	strings.Replace(reqData.Location, " ", "+", -1)
	strings.Replace(reqData.Title, "-", "+", -1)
	strings.Replace(reqData.Location, "-", "+", -1)

	//instantiate the collectors

	//INDEED
	cIndeed := colly.NewCollector()

	scrapeIndeed(cIndeed)

	// STACK EXCHANGE JOBS
	cStackOverFlow := colly.NewCollector()
	scrapeStackOverFlow(cStackOverFlow)

	cStackOverFlow.Visit("https://stackoverflow.com/jobs?q=" + reqData.Title + "&l=" + reqData.Location)

	cIndeed.Visit("https://www.indeed.com/jobs?q=" + reqData.Title + "&l=" + reqData.Location + "&explvl=entry_level")

	//encode the map string into json

	//replace explicit unicode with &
	jobsJson, _ := json.Marshal(jobs)
	jobsJson = bytes.Replace(jobsJson, []byte("\\u0026"), []byte("&"), -1)

	// Set responsewriter's header to let client expect json
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Fprintf(w, string(jobsJson))

	//clear array for next request
	jobs = make([]Job, 0, 300)
	r.Body.Close()

}

func main() {
	fmt.Println("Visiting:")
	http.HandleFunc("/api/", Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
