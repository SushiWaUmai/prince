package utilcmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func PermissionCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) < 1 {
		return utils.CreateTextMessage("Usage: permission <permission> <user>"), errors.New("Not enough arguments")
	}

	perm := strings.ToUpper(args[0])

	// NONE, USER, ADMIN
	if perm != "NONE" && perm != "USER" && perm != "ADMIN" {
		return utils.CreateTextMessage("Invalid permission type. Available: NONE, USER, ADMIN"), errors.New("Invalid permission type")
	}

	isGroup := chat.Server == "g.us"
	if isGroup && (ctx == nil || len(ctx.MentionedJid) <= 0) {
		return utils.CreateTextMessage("No user mentioned"), errors.New("No user mentioned")
	}

	if isGroup {
		for _, u := range ctx.MentionedJid {
			jid, err := types.ParseJID(u)
			if err != nil {
				return nil, errors.New("Failed to parse JID")
			}

			err = db.UpdateUserPermission(jid.ToNonAD().User, perm)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err := db.UpdateUserPermission(chat.ToNonAD().User, perm)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("permission", "OP", PermissionCommand)
}
