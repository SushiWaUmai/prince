package utilcmds

import (
	"errors"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func UnaliasCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) < 1 {
		return nil, errors.New("No name provided")
	}
	name := args[0]

	err := db.DeleteAlias(name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("unalias", "OP", UnaliasCommand)
}
