package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"./util"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
var api = slack.New(os.Getenv("SLACK_BOT_TOKEN"))

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
				fmt.Println("--------------------")
				if len(ev.Attachments) > 0 {
					for i := 0; i < len(ev.Attachments); i++ {
						fmt.Println(ev.Attachments[i].Text)
						if strings.Contains(ev.Attachments[i].Text, "*Finish*: No → Yes") {
							client := util.NotionClient("70cf786c353be577c943b2ccbde5b224dc7170e5c910b1ed96176c24cf3122330a5a0d3dc998f5008f239c839854cab6a483ae8653ffb6426b6cda5925601bc7f5471662f4bd9285fed87ebd3a22")
							util.NotionExport(client, 2)
							fmt.Println("Got it")
						}
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
