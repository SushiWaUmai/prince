package textcmds

import (
	"github.com/SushiWaUmai/prince/utils"
	openai "github.com/sashabaranov/go-openai"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func init() {
	utils.CreateCommand("clearchat", "ADMIN", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		PastMessages[chat] = make([]openai.ChatCompletionMessage, 0)

		return nil, nil
	})
}
