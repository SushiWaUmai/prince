package utilcmds

import (
	"github.com/SushiWaUmai/prince/utils"
	"github.com/yuin/gopher-lua"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func LuaCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	L := lua.NewState()
	defer L.Close()

	err := L.DoString(`print("hello")`); 


	if err != nil {
		return nil, err
	}

	return nil, nil
}

func init() {
	utils.CreateCommand("lua", "ADMIN", LuaCommand)
}
