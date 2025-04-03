package channels

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/OssiPesonen/builder-app-go/internal/services/broadcaster"
)

type BlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Block struct {
	Type string    `json:"type"`
	Text BlockText `json:"text"`
}

type SlackMessage struct {
	Text   string  `json:"text"`
	Blocks []Block `json:"blocks"`
}

type Slack struct {
	URL     string
	Channel <-chan broadcaster.Message
}

func NewSlack(url string, ch <-chan broadcaster.Message) broadcaster.Subscriber {
	return &Slack{
		URL:     url,
		Channel: ch,
	}
}

func (s *Slack) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range s.Channel {
		s.Publish(msg.Title, msg.Body)
	}
}

func (s *Slack) Publish(title string, markdown string) {
	slackMessage := SlackMessage{Text: title, Blocks: []Block{
		{
			Type: "section",
			Text: BlockText{
				Type: "mrkdwn",
				Text: markdown,
			},
		},
	}}

	payload, err := json.Marshal(slackMessage)
	if err != nil {
		return
	}

	r, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return
	}

	defer res.Body.Close()
}
