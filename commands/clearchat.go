package commands

import (
	"github.com/SushiWaUmai/prince/utils"
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

func init() {
	utils.CreateCommand("clearchat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		PastMessages = make([]openai.ChatCompletionMessage, 0)

		return nil, nil
	})
}
