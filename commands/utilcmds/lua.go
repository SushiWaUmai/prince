package utilcmds

import (
	"context"
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/yuin/gopher-lua"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"layeh.com/gopher-luar"
)

func LuaCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var script string
	if len(args) > 0 {
		script = strings.Join(args, " ")
	} else {
		script, _ = utils.GetTextContext(pipe)
	}
	script = strings.TrimSpace(script)

	if script == "" {
		return utils.CreateTextMessage("Please input a lua script"), errors.New("No script specified")
	}

	sendMessage := func(text string) {
		client.SendMessage(context.Background(), chat, utils.CreateTextMessage(text))
	}

	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("sendMessage", luar.New(L, sendMessage))
	err := L.DoString(script)

	if err != nil {
		return utils.CreateTextMessage(err.Error()), err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("lua", "ADMIN", LuaCommand)
}
