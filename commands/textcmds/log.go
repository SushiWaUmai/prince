package textcmds

import (
	"errors"
	"log"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func LogCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) > 0 {
		text := strings.Join(args, " ")
		text = strings.TrimSpace(text)
		log.Println(text)

		return pipe, nil
	}

	if pipe == nil {
		return nil, errors.New("No pipe provided")
	}

	text, _ := utils.GetTextContext(pipe)
	log.Println(text)

	return pipe, nil
}

func init() {
	utils.CreateCommand("log", "OP", "Logs the provided text or replied message", LogCommand)
}
