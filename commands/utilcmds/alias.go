package utilcmds

import (
	"errors"
	"log"
	"strings"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func aliasCreate(args []string, pipe *waProto.Message) (*waProto.Message, error) {
	if len(args) < 1 {
		return nil, errors.New("No name provided")
	}
	name := args[0]
	log.Println(name)

	var content string
	if len(args) > 1 {
		content = strings.Join(args[1:], " ")
	} else {
		content, _ = utils.GetTextContext(pipe)
	}
	content = strings.TrimSpace(content)

	if content == "" {
		return nil, errors.New("No content provided")
	}

	db.UpsertAlias(name, content)

	return nil, nil
}

func aliasDelete(args []string) (*waProto.Message, error) {
	if len(args) < 1 {
		return nil, errors.New("No name provided")
	}
	name := args[0]

	db.DeleteAlias(name)
	return nil, nil
}

func aliasGet(args []string) (*waProto.Message, error) {
	if len(args) < 1 {
		return nil, errors.New("No name provided")
	}
	name := args[0]

	alias := db.GetAlias(name)
	if alias == nil {
		response := &waProto.Message{
			Conversation: proto.String("Alias with name \"" + name + "\" not found"),
		}
		return response, nil
	}

	response := &waProto.Message{
		Conversation: proto.String(alias.Name + ": " + alias.Content),
	}

	return response, nil
}

func AliasCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) < 2 {
		return nil, errors.New("Not enough arguments")
	}

	op := strings.ToLower(args[0])
	switch op {
	case "create":
		return aliasCreate(args[1:], pipe)
	case "delete":
		return aliasDelete(args[1:])
	case "get":
		return aliasGet(args[1:])
	default:
		return nil, errors.New("Invalid operation")
	}
}

func init() {
	utils.CreateCommand("alias", "OP", AliasCommand)
}
