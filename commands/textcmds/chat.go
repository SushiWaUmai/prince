package textcmds

import (
	"context"
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

// TODO: Move this to database
var PastMessages = make(map[types.JID][]openai.ChatCompletionMessage, 0)

func init() {
	utils.CreateCommand("chat", "ADMIN", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var prompt string

		pipeString, _ := utils.GetTextContext(pipe)
		if pipeString != "" {
			prompt = pipeString + "\n\n"
		}

		if len(args) > 0 {
			prompt += strings.Join(args, " ")
		}

		prompt = strings.TrimSpace(prompt)

		if len(prompt) <= 0 {
			return nil, errors.New("Failed to generate openai response, no prompt was provided")
		}

		PastMessages[chat] = append(PastMessages[chat], openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		})

		resp, err := utils.OpenAIClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: PastMessages[chat],
			},
		)

		if err != nil {
			return nil, err
		}

		reply := resp.Choices[0].Message.Content

		PastMessages[chat] = append(PastMessages[chat], openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: reply,
		})

		response := &waProto.Message{
			Conversation: &reply,
		}

		return response, nil
	})
}
