package utilcmds

import (
	"context"
	"errors"
	"strings"
	"time"

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

	msgSent := 0
	sendMessage := func(text string) {
		if msgSent < 16 {
			client.SendMessage(context.Background(), chat, utils.CreateTextMessage(text))
			msgSent++
		}
	}
	getArg := func(i int) string {
		if i >= len(luaArgs) {
			return ""
		}

		return luaArgs[i]
	}

	L := lua.NewState()
	defer L.Close()

	luaCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	L.SetContext(luaCtx)

	L.SetGlobal("sendMessage", luar.New(L, sendMessage))
	L.SetGlobal("getArg", luar.New(L, getArg))
	err := L.DoString(script)

	if err != nil {
		return utils.CreateTextMessage(err.Error()), err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("lua", "ADMIN", LuaCommand)
}
