package mediacmds

import (
	"errors"
	"strconv"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var voicevox = utils.CreateVoiceVoxClient()

func VoicevoxCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if len(args) <= 0 {
		return nil, errors.New("Please enter a valid speaker index")			
	}

	var text string
	speakerIndex, err := strconv.Atoi(args[0])
	if len(args) > 1 {
		text = strings.Join(args[1:], " ")
	} else {
		text, _ = utils.GetTextContext(pipe)
	}
	text = strings.TrimSpace(text)

	if text == "" {
		return utils.CreateTextMessage("Please specify a text to speak"), errors.New("No Text specified")
	}

	query, err := voicevox.CreateQuery(speakerIndex, text)

	if err != nil {
		return nil, err
	}

	buffer, err := voicevox.CreateVoice(speakerIndex, true, query)
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateAudioMessage(client, buffer)

	return response, nil
}

func init() {
	utils.CreateCommand("voicevox", "USER", "Converts text to speech using specified speaker index and message.", VoicevoxCommand)
}
