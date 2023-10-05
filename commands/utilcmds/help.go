package utilcmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func HelpCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	// if an argument is given try to find the the command
	// the command is not given return an error message
	// if it is return a message with the permission and help message
	if len(args) > 0 {
		cmd, ok := utils.CommandMap[args[0]]
		if !ok {
			return utils.CreateTextMessage("Command \"" + args[0] + "\" not found"), nil
		}

		var data []string
		data = append(data, "*"+cmd.Name+"*")
		data = append(data, cmd.Permission)
		data = append(data, cmd.Description)

		return utils.CreateTextMessage(strings.Join(data, "\n")), nil
	}

	var cmds []string

	for _, c := range utils.CommandMap {
		cmds = append(cmds, c.Name)
	}

	return utils.CreateTextMessage(strings.Join(cmds, "\n")), nil
}

func init() {
	utils.CreateCommand("help", "USER", "Sends a list of commands", HelpCommand)
}
