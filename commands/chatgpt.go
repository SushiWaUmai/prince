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

var OpenAIClient = openai.NewClient(env.OPENAI_API_KEY)

func init() {
	createCommand("chatgpt", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) {
		var prompt string

		if ctx != nil && ctx.QuotedMessage != nil {
			if ctx.QuotedMessage.Conversation != nil {
				prompt = *ctx.QuotedMessage.Conversation + "\n\n"
			} else if ctx.QuotedMessage.AudioMessage != nil {
				// Download the voice message
				audioData, err := client.Download(ctx.QuotedMessage.AudioMessage)
				if err != nil {
					client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
						Conversation: proto.String("Failed to download the voice message"),
					})
					return
				}

				text, err := TranscribeAudio(audioData)
				if err != nil {
					client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
						Conversation: proto.String("Failed to download the voice message"),
					})
					return
				}

				prompt = text + "\n\n"
			}
		}

		if len(args) > 0 {
			prompt += strings.Join(args, " ")
		}

		prompt = strings.TrimSpace(prompt)

		if len(prompt) <= 0 {
			log.Println("Failed to generate openai response, no prompt was provided")
			return
		}

		resp, err := OpenAIClient.CreateChatCompletion(
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

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})

		if err != nil {
			log.Println(err)
		}
	})
}
