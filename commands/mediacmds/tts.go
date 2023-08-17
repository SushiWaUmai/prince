package mediacmds

import (
	"errors"

	"os"
	"strings"

	"github.com/SushiWaUmai/prince/utils"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

func TextToSpeechCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
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

	speech := htgotts.Speech{
		Folder:   ".",
		Language: voices.English,
	}

	_, err := speech.CreateSpeechFile(text, "speach")
	if err != nil {
		return nil, err
	}
	defer os.Remove("speach.mp3")

	// get the bytes
	audioBytes, err := os.ReadFile("speach.mp3")
	if err != nil {
		return nil, err
	}

	response, err := utils.CreateAudioMessage(client, audioBytes)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func init() {
	utils.CreateCommand("tts", "USER", TextToSpeechCommand)
}
