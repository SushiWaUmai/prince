package commands

import (
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	createCommand("clearchat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) error {
		PastMessages = make([]openai.ChatCompletionMessage, 0)

		return nil
	})
}
