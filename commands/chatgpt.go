package commands

import (
	"context"
	"log"
	"strings"

	"github.com/SushiWaUmai/prince/env"
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var chatgpt *openai.Client

func init() {
	chatgpt = openai.NewClient(env.OPENAI_API_KEY)

	createCommand("chatgpt", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		var prompt string

		if len(args) > 0 {
			prompt = strings.Join(args, " ")
		} else if ctx.QuotedMessage != nil {
			prompt = *ctx.QuotedMessage.Conversation
		} else {
			log.Println("Failed to generate openai response, no prompt was provided")
			return
		}

		resp, err := chatgpt.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)

		if err != nil {
			log.Printf("ChatCompletion error: %v\n", err)
			return
		}

		reply := resp.Choices[0].Message.Content

		client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: proto.String(reply),
		})
	})
}
