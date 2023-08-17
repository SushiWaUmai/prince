package client

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var maxDepth = 10

func appendCommandRecursive(cmd string, commandInput []utils.CommandInput, i int) ([]utils.CommandInput, error) {
	if i > maxDepth {
		return nil, errors.New("Max Depth reached")
	}

	// split the command name with the arguments
	cmd = strings.TrimSpace(cmd)
	split := strings.Split(cmd, " ")

	cmdName := strings.ToLower(split[0])
	cmdArgs := split[1:]

	_, ok := utils.CommandMap[cmdName]
	if !ok {
		// get alias
		alias, err := db.GetAlias(cmdName)
		if err != nil {
			return commandInput, err
		}

		// TODO: Make it work with AST
		commandsSplit := strings.Split(alias.Content, "|")

		for _, c := range commandsSplit {
			var err error
			commandInput, err = appendCommandRecursive(c, commandInput, i+1)

			if err != nil {
				return nil, err
			}
		}

		aliasArgs := strings.Split(alias.Content, " ")
		cmdName = aliasArgs[0]
		if len(aliasArgs) > 1 {
			cmdArgs = append(aliasArgs[1:], cmdArgs...)
		}
	}

	return append(commandInput,
		utils.CommandInput{
			Name: cmdName,
			Args: cmdArgs,
		},
	), nil
}

func appendCommand(cmd string, commandInput []utils.CommandInput) ([]utils.CommandInput, error) {
	return appendCommandRecursive(cmd, commandInput, 0)
}

func (client *PrinceClient) handleCommand(message *waProto.Message, msgId types.MessageID, chat types.JID, user string) {
	content, ctx := utils.GetTextContext(message)

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	perm := db.GetUserPermission(user).Permission
	fromMe := client.wac.Store.ID.ToNonAD().User == user
	if fromMe {
		perm = "OP"
	}

	if perm == "NONE" {
		return
	}

	content = content[len(client.commandPrefix):]
	content = strings.TrimSpace(content)

	commandsSplit := strings.Split(content, "|")
	var commandInput []utils.CommandInput

	for _, c := range commandsSplit {
		var err error
		commandInput, err = appendCommand(c, commandInput)

		if err != nil {
			log.Println(err)
			return
		}
	}

	// Validate all commands
	for _, c := range commandInput {
		cmd, ok := utils.CommandMap[c.Name]
		if !ok {
			log.Println("Command not found: ", c.Name)
			return
		}
		if !db.ComparePermission(perm, cmd.Permission) {
			log.Println("Not enough permission")
			return
		}
	}

	var pipe *waProto.Message = nil
	if ctx != nil {
		pipe = ctx.QuotedMessage
	}

	reaction := "⏳"

	client.SendMessage(chat, &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid: proto.String(chat.String()),
				FromMe:    &fromMe,
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	var err error
	for _, c := range commandInput {
		log.Println("Runnning commmand", c.Name, "with args", c.Args)
		pipe, err = utils.CommandMap[c.Name].Execute(client.wac, chat, user, ctx, pipe, c.Args)

		if err == nil {
			reaction = "👍"
		} else {
			log.Println(err)
			reaction = "❌"
			break
		}
	}

	client.SendMessage(chat, &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid: proto.String(chat.String()),
				FromMe:    &fromMe,
				Id:        &msgId,
			},
			Text:              &reaction,
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})

	if pipe != nil {
		client.SendCommandMessage(chat, user, pipe)
	}

	log.Println("Done.")
}
