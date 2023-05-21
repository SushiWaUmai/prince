package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("help", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		var cmds []string

		for _, c := range utils.CommandMap {
			cmds = append(cmds, c.Name)
		}

		response := &waProto.Message{
			Conversation: proto.String(strings.Join(cmds, "\n")),
		}
		return response, nil
	})
}
