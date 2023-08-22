package mediacmds

import (
	"errors"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var voicevox = utils.CreateVoiceVoxClient()
var zundamonIdx = 1

// TODO: Change this to voicevox command and add zundamon alias
func ZundamonCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	var text string
	if len(args) > 0 {
		text = strings.Join(args, " ")
	} else {
		text, _ = utils.GetTextContext(pipe)
	}
	text = strings.TrimSpace(text)

	if text == "" {
		return utils.CreateTextMessage("Please specify a text to speak"), errors.New("No Text specified")
	}

	query, err := voicevox.CreateQuery(zundamonIdx, text)

	if err != nil {
		return nil, err
	}

	buffer, err := voicevox.CreateVoice(1, true, query)
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateAudioMessage(client, buffer)

	return response, nil
}

func init() {
	utils.CreateCommand("zundamon", "USER", "", ZundamonCommand)
}
