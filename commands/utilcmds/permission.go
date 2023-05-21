package utilcmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func init() {
	utils.CreateCommand("permission", "OP", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		if len(args) < 1 {
			response := &waProto.Message{
				Conversation: proto.String("Usage: permission <permission> <user>"),
			}
			return response, errors.New("Not enough arguments")
		}

		perm := strings.ToUpper(args[0])

		// NONE, USER, ADMIN
		if perm != "NONE" && perm != "USER" && perm != "ADMIN" {
			response := &waProto.Message{
				Conversation: proto.String("Invalid permission type. Available: NONE, USER, ADMIN"),
			}
			return response, errors.New("Invalid permission type")
		}

		isGroup := chat.Server == "g.us"
		if isGroup && (ctx == nil || len(ctx.MentionedJid) <= 0) {
			response := &waProto.Message{
				Conversation: proto.String("No user mentioned"),
			}
			return response, errors.New("No user mentioned")
		}

		if isGroup {
			for _, u := range ctx.MentionedJid {
				jid, err := types.ParseJID(u)
				if err != nil {
					return nil, errors.New("Failed to parse JID")
				}

				db.UpdateUserPermission(jid.ToNonAD().User, perm)
			}
		} else {
			db.UpdateUserPermission(chat.ToNonAD().User, perm)
		}

		response := &waProto.Message{
			Conversation: proto.String("Permission set"),
		}
		return response, nil
	})
}
