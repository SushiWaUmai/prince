package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func HelpCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var cmds []string

	for _, c := range utils.CommandMap {
		cmds = append(cmds, c.Name)
	}

	return utils.CreateTextMessage(strings.Join(cmds, "\n")), nil
}

func init() {
	utils.CreateCommand("help", "USER", HelpCommand)
}
