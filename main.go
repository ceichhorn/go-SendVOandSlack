package main

import (
	"./client"
        "github.com/bluele/slack"
        "os"
	"log"
        "fmt"
	"time"
 "net/http"
 "io/ioutil"
 "encoding/json"
)

type Data struct {
   Data string
   Sha string `json:"sha"`
   URL string `json:"html_url"`
    Commit Commit `json:"commit"`
}

type Commit struct {
    Author Author `json:"author"`
    Committer Committer `json:"committer"`
}
type Committer  struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Date string `json:"date"`
}

type Author struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Date string `json:"date"`
}

const (
    USERNAME = "GITHUB_USER"
    PASSWORD = "GIT_HUB_TOKEN"
    URL      = "GIT_HUB_URL"
//    URL      = "https://api.github.com/repos/GROUP/REPO/commits?page=0"
    channelName = "test-channel"
)

var Cemail string
var Cuser string

func get_content() {
    req, err := http.NewRequest("GET", URL, nil)
    req.SetBasicAuth(USERNAME, PASSWORD)

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")

    cli := &http.Client{}
    resp, err := cli.Do(req)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var data []Data
    err = json.Unmarshal(body, &data)
    if err != nil {
        panic(err.Error())
    }


  for i := 0; i < len(data); i++ {
     C1 := "GitHub"
       if (data[i].Commit.Committer.Name) != C1 {
	fmt.Println("Committer: " + data[i].Commit.Committer.Name)
	fmt.Println("Email: " + data[i].Commit.Committer.Email)
        Cemail = (data[i].Commit.Committer.Email)
        Cuser = (data[i].Commit.Committer.Name)
        fmt.Println (Cemail)
        break
        }
       }
 }

////

func post_voalert() {
VO_KEY := (os.Getenv("VO_ROUTE_KEY"))
Msg := fmt.Sprintf("There is an issue with the X. As this may be an issue with config or code, we are now looping in additional resources to troubleshoot at those levels. We are requesting that %s, the last author of a merged config commit assist.  Troubleshooting is happening in the #Troubleshooting Slack channel", Cuser)
	vo := victorops.NewClient(VO_KEY)

	e := &victorops.Event{
		RoutingKey:        "VO_GROUP_KEY",
		MessageType:       victorops.Critical,
		EntityID:          "",
		StateMessage:      Msg,
		Timestamp:         time.Now(),
		EntityDisplayName: "",
	}

        fmt.Println(vo)
        fmt.Println(e)
        fmt.Println(Cuser)
        fmt.Println(Cemail)
	resp, err := vo.SendAlert(e)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	log.Printf("Response: %+v\n", resp)
}

//////

func post_slack() {
Msg := fmt.Sprintf("There is an issue with the X. As this may be an issue with config or code, we are now looping in additional resources to troubleshoot at those levels. We are requesting that %s, the last author of a merged config commit assist.  Troubleshooting is happening in the #Troubleshooting Slack channel", Cuser)

SLACK := (os.Getenv("SLACK_TOKEN"))
api := slack.New(SLACK)
  err := api.ChatPostMessage(channelName, Msg, nil)
  if err != nil {
    panic(err)
  }
}
func main() {
    get_content()
  fmt.Println("Main email:" + Cemail)
  fmt.Println(Cuser)
    post_voalert()
    post_slack()

}
