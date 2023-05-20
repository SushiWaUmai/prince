package commands

import (
	"context"
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// TODO: Move this to database
var PastMessages = make([]openai.ChatCompletionMessage, 0)

func init() {
	createCommand("chat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var prompt string

		pipeString, _ := GetTextContext(pipe)
		if pipeString != "" {
			prompt = pipeString + "\n\n"
		}

		if pipe != nil && pipe.AudioMessage != nil {
			// Download the voice message
			audioData, err := client.Download(pipe.AudioMessage)
			if err != nil {
				client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
					Conversation: proto.String("Failed to download voice message"),
				})
				return nil, err
			}

			text, err := utils.TranscribeAudio(audioData)
			if err != nil {
				client.SendMessage(context.Background(), messageEvent.Info.Chat, &waProto.Message{
					Conversation: proto.String("Failed to transcribe voice message"),
				})
				return nil, err
			}

			prompt = text + "\n\n"
		}

		if len(args) > 0 {
			prompt += strings.Join(args, " ")
		}

		prompt = strings.TrimSpace(prompt)

		if len(prompt) <= 0 {
			return nil, errors.New("Failed to generate openai response, no prompt was provided")
		}

		PastMessages = append(PastMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		})

		resp, err := utils.OpenAIClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: PastMessages,
			},
		)

		if err != nil {
			return nil, err
		}

		reply := resp.Choices[0].Message.Content

		PastMessages = append(PastMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})

		response := &waProto.Message{
			Conversation: &reply,
		}

		return response, nil
	})
}
