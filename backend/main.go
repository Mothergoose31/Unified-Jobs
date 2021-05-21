package main

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
type GithubJson []struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	CompanyURL  string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HowToApply  string `json:"how_to_apply"`
	CompanyLogo string `json:"company_logo"`
}

var reqData = RequestData{Title: "", Location: ""}
var jobs = make([]Job, 0, 300)

func main() {

}
