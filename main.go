package main

import (
    //"strings"
	"bytes"
	"encoding/json"
	"fmt"
    "os"
	"net/http"
	"regexp"
	"github.com/tamama9527/notion_bot/util"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))
var client = util.NotionClient("c3e704943239cc91d202a0e97b58ceaf6e219e5ccef7685fb5b863026466d766abd2cc220f21e532f43579b23bc2da4e5e3c669ab3b63a1b224d068cff82176629047c8fa0b6cad8d71135344ae6")
func main() {
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_EVENT_TOKEN")}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
            case *slackevents.MessageEvent:
				fmt.Println(ev.Text)
				str := []byte(ev.Text)
				reg := regexp.MustCompile(`\*<(\S*)\|(.*)>\*`)
				result := reg.FindSubmatch(str)
				url := string(result[1])
				title := string(result[2])
                regPageid := regexp.MustCompile(`-(\w+)\?|\/(\w+)\?`)
                TempResult := regPageid.FindStringSubmatch(url)
				PageID := ""
				if len(TempResult[1]) > len(TempResult[2]){
					PageID = TempResult[1]
					fmt.Println(PageID)
				} else{
					PageID = TempResult[2]
					fmt.Println(PageID)
				}
				fmt.Println(PageID)
                fmt.Println("--------------------")
                fmt.Println(url)
                fmt.Println(title)

                if len(ev.Attachments) > 0{
                    for i:= 0; i < len(ev.Attachments); i++{
                        fmt.Println(ev.Attachments[i].Text)
                        // if strings.Contains(ev.Attachments[i].Text,"*Finish*: No → Yes"){
						// 	util.NotionExport(client,PageID)
						// 	fmt.Println("Got it")
                        // }
						util.NotionExport(client,PageID)
                        util.NotionPages(client,PageID)
						fmt.Println("Got it")
                    }
                } else {
                    fmt.Println(ev.Text)
                }
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	fmt.Println(os.Getenv("SLACK_EVENT_TOKEN"))
	fmt.Println(os.Getenv("SLACK_BOT_TOKEN"))
	http.ListenAndServe(":3000", nil)
}
