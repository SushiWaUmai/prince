package utilcmds

import (
	"errors"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func AliasCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) < 1 {
		return nil, errors.New("No name provided")
	}
	name := args[0]

	var content string
	if len(args) > 1 {
		content = strings.Join(args[1:], " ")
	} else {
		content, _ = utils.GetTextContext(pipe)
	}
	content = strings.TrimSpace(content)

	if content == "" {
		alias := db.GetAlias(name)
		if alias == nil {
			return utils.CreateTextMessage("Alias with name \"" + name + "\" not found"), nil
		}

		return utils.CreateTextMessage(alias.Name + "=\"" + alias.Content + "\""), nil
	}

	err := db.UpsertAlias(name, content)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("alias", "OP", AliasCommand)
}
