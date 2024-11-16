package client

import (
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type Line struct {
	Bot         *messaging_api.MessagingApiAPI
	To          string
	BaseMessage string
}

func NewLine(token string, to string, baseMessage string) *Line {
	bot, err := messaging_api.NewMessagingApiAPI(token)
	if err != nil {
		panic(err)
	}
	return &Line{
		Bot:         bot,
		To:          to,
		BaseMessage: baseMessage,
	}
}

func (l *Line) Notify(message string) {
	if l.To == "" {
		l.Bot.Broadcast(&messaging_api.BroadcastRequest{
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message,
				},
			},
		}, "")
	} else {
		l.Bot.PushMessage(&messaging_api.PushMessageRequest{
			To: l.To,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message,
				},
			},
		}, "")
	}
}

func (l *Line) CreateMessage(url string) string {
	return strings.ReplaceAll(l.BaseMessage, "{url}", url)
}
