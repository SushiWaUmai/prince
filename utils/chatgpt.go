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

func SetSystemChat(chat types.JID, content string) {
	// check if there is any openai.ChatMessageRoleSystem
	// if there is, override
	for i, msg := range PastMessages[chat] {
		if msg.Role == openai.ChatMessageRoleSystem {
			PastMessages[chat][i].Content = content
			return
		}
	}

	// if there isn't, add at the begining
	PastMessages[chat] = append([]openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: content,
		},
	}, PastMessages[chat]...)
}
