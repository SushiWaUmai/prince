package utilcmds

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	princeClient "github.com/SushiWaUmai/prince/client"
	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/utils"
	"github.com/yuin/gopher-lua"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"layeh.com/gopher-luar"
)

func LuaCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var script string
	var luaArgs []string
	if pipe == nil {
		script = args[0]
		luaArgs = args[1:]
	} else {
		script, _ = utils.GetTextContext(pipe)
		luaArgs = args
	}
	script = strings.TrimSpace(script)

	if script == "" {
		return utils.CreateTextMessage("Please input a lua script"), errors.New("No script specified")
	}

	var err error
	msgSent := 0
	commandExecuted := 0
	resultMessage := ""

	sendMessage := func(text string) {
		if msgSent < 16 {
			client.SendMessage(context.Background(), chat, utils.CreateTextMessage(text))
			msgSent++
		} else {
			err = errors.New("Cannot send more than 16 messages from single a lua script")
		}
	}

	executeCommand := func(text string) {
		if commandExecuted < 8 {
			var message *waProto.Message
			message, err = princeClient.RunCommand(client, string(env.BOT_PREFIX)+text, ctx, chat, user)
			if err != nil {
				return
			}

			client.SendMessage(context.Background(), chat, message)
			commandExecuted++
		} else {
			err = errors.New("Cannot execute more than 8 command from a single lua script")
		}
	}

	printMessage := func(data interface{}) {
		log.Println(data)
		if str, ok := data.(string); ok {
			resultMessage += str + "\n"
		}
	}

	L := lua.NewState()
	defer L.Close()

	luaCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	L.SetContext(luaCtx)

	L.SetGlobal("sendMessage", luar.New(L, sendMessage))
	L.SetGlobal("executeCommand", luar.New(L, executeCommand))
	L.SetGlobal("print", luar.New(L, printMessage))

	for _, command := range utils.CommandMap {
		cmd := command
		execute := func(arg ...string) string {
			msg, err := princeClient.RunCommand(client, string(env.BOT_PREFIX)+cmd.Name+" "+strings.Join(arg, " "), ctx, chat, user)
			if err != nil {
				log.Println(err)
				return ""
			}
			if msg.Conversation == nil {
				return ""
			}

			return *msg.Conversation
		}

		L.SetGlobal(cmd.Name, luar.New(L, execute))
	}

	lTable := L.NewTable()
	for i, str := range luaArgs {
		L.RawSetInt(lTable, i+1, lua.LString(str))
	}

	L.SetGlobal("arg", lTable)

	luaErr := L.DoString(script)

	if luaErr != nil {
		return utils.CreateTextMessage(luaErr.Error()), luaErr
	}

	if err != nil {
		return utils.CreateTextMessage(err.Error()), err
	}

	if resultMessage == "" {
		return nil, nil
	}

	return utils.CreateTextMessage("```" + resultMessage + "```"), nil
}

func init() {
	utils.CreateCommand("lua", "ADMIN", "Executes a lua script", LuaCommand)
}
