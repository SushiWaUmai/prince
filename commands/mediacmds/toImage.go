package mediacmds

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/SushiWaUmai/prince/utils"
	"github.com/chai2010/webp"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func ToImageCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	if pipe == nil || pipe.StickerMessage == nil {
		response := &waProto.Message{
			Conversation: proto.String("Please reply to a sticker message"),
		}
		return response, errors.New("No StickerMessage quoted")
	}
	stickerImg := pipe.StickerMessage

	buffer, err := client.Download(stickerImg)
	if buffer == nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}

	webpByte, err := webp.EncodeRGBA(img, *proto.Float32(1))
	if err != nil {
		return nil, err
	}

	imgMsg, err := utils.CreateImgMessage(client, webpByte)
	if err != nil {
		return nil, err
	}

	response := &waProto.Message{
		ImageMessage: imgMsg,
	}
	return response, nil
}

func init() {
	utils.CreateCommand("toimage", "USER", ToImageCommand)
}
