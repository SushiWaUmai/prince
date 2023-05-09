package commands

import (
	"context"
	"errors"
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
	createCommand("chatgpt", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, args []string) error {
		var prompt string

		if ctx != nil && ctx.QuotedMessage != nil {
			if ctx.QuotedMessage.Conversation != nil {
				prompt = *ctx.QuotedMessage.Conversation + "\n\n"
			} else if ctx.QuotedMessage.AudioMessage != nil {
				// Download the voice message
				audioData, err := client.Download(ctx.QuotedMessage.AudioMessage)
				if err != nil {
					client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
						Conversation: proto.String("Failed to download voice message"),
					})
					return err
				}

				text, err := TranscribeAudio(audioData)
				if err != nil {
					client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
						Conversation: proto.String("Failed to transcribe voice message"),
					})
					return err
				}

				prompt = text + "\n\n"
			}
		}

		if len(args) > 0 {
			prompt += strings.Join(args, " ")
		}

		prompt = strings.TrimSpace(prompt)

		if len(prompt) <= 0 {
			return errors.New("Failed to generate openai response, no prompt was provided")
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
			return err
		}

		reply := resp.Choices[0].Message.Content

		_, err = client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
			Conversation: &reply,
		})

		if err != nil {
			return err
		}

		return nil
	})
}
