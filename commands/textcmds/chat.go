package textcmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func ChatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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

	reply, err := utils.GetChatReponse(chat, prompt)
	if err != nil {
		return nil, err
	}

	return utils.CreateTextMessage(reply), nil
}

func init() {
	utils.CreateCommand("chat", "ADMIN", "Generates a chat response based on the provided or replied text using ChatGPT", ChatCommand)
}
