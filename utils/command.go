package utils

import (
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

var CommandMap map[string]Command = make(map[string]Command)

type CommandInput struct {
	Name string
	Args []string
}

type Command struct {
	Name    string
	Execute func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error)
}

func CreateCommand(name string, execute func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error)) {
	CommandMap[name] = Command{
		Name:    name,
		Execute: execute,
	}
}

func GetTextContext(msg *waProto.Message) (string, *waProto.ContextInfo) {
	if msg == nil {
		return "", nil
	}

	var text string = ""
	var ctx *waProto.ContextInfo

	if msg_type := msg.GetAudioMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
	} else if msg_type := msg.GetButtonsMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetContentText()
	} else if msg_type := msg.GetButtonsResponseMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetSelectedButtonId()
	} else if msg_type := msg.GetContactMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetVcard()
	} else if msg_type := msg.GetConversation(); msg_type != "" {
		text = msg_type
	} else if msg_type := msg.GetDocumentMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetFileName()
	} else if msg_type := msg.GetExtendedTextMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetText()
	} else if msg_type := msg.GetImageMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetCaption()
	} else if msg_type := msg.GetListMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetDescription()
	} else if msg_type := msg.GetListResponseMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.SingleSelectReply.GetSelectedRowId()
	} else if msg_type := msg.GetProductMessage(); msg_type != nil {
		ctx = msg_type.ContextInfo
		text = msg_type.GetBody()
	} else if msg_type := msg.GetReactionMessage(); msg_type != nil {
		text = msg_type.GetText()
	} else if msg_type := msg.GetStickerMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
	} else if msg_type := msg.GetTemplateButtonReplyMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetSelectedId()
	} else if msg_type := msg.GetTemplateMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		if msg_subtype := msg_type.GetFourRowTemplate(); msg_subtype != nil {
			text = msg_subtype.Content.GetNamespace()
		} else if msg_subtype := msg_type.GetHydratedTemplate(); msg_subtype != nil {
			text = msg_subtype.GetTemplateId()
		} else if msg_subtype := msg_type.GetHydratedFourRowTemplate(); msg_subtype != nil {
			text = msg_subtype.GetTemplateId()
		}
	} else if msg_type := msg.GetVideoMessage(); msg_type != nil {
		ctx = msg_type.GetContextInfo()
		text = msg_type.GetCaption()
	}

	return text, ctx
}
