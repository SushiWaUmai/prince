package utils

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow/types"
)

// TODO: Move this to database
var PastMessages = make(map[types.JID][]openai.ChatCompletionMessage, 0)

func GetChatReponse(chat types.JID, prompt string) (string, error) {
	PastMessages[chat] = append(PastMessages[chat], openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})

	resp, err := OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: PastMessages[chat],
		},
	)

	if err != nil {
		return "", err
	}

	reply := resp.Choices[0].Message.Content

	PastMessages[chat] = append(PastMessages[chat], openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply,
	})

	return reply, nil
}

func ClearChat(chat types.JID) {
	PastMessages[chat] = make([]openai.ChatCompletionMessage, 0)
}
