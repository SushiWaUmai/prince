package mediacmds

import (
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func StableDiffusionCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var prompt string
	pipeString, _ := utils.GetTextContext(pipe)
	if pipeString != "" {
		prompt = pipeString
	}

	if len(args) > 0 {
		if prompt != "" {
			prompt += ", "
		}

		prompt += strings.Join(args, " ")
	}

	buffer, err := utils.Txt2Img(prompt)
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateImgMessage(client, buffer)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func init() {
	utils.CreateCommand("sd", "USER", "Transforms user-provided text into an image using AI", StableDiffusionCommand)
}
