package client

import (
	"log"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"mvdan.cc/xurls/v2"
)

func (client *PrinceClient) handleMessage(e *events.Message) error {
	client.handleCommand(e.Message, e.Info.ID, e.Info.Chat, e.Info.Sender.User)

	err := client.handleMessageEvents(e)
	if err != nil {
		return err
	}

	return err
}

func (client *PrinceClient) handleMessageEvents(e *events.Message) error {
	jid := e.Info.Chat.String()

	msgEvents, err := db.GetMessageEvents(jid)
	if err != nil {
		return err
	}

	for _, evt := range msgEvents {
		switch evt.Type {
		case "DOWNLOAD":
			client.handleMessageDownload(e)
		case "CHAT":
			client.handleMessageChat(e)
		}
	}

	return nil
}

func (client *PrinceClient) handleMessageDownload(e *events.Message) {
	content, _ := utils.GetTextContext(e.Message)

	rxStrict := xurls.Strict()
	fetchUrl := rxStrict.FindString(content)

	msg, err := utils.GetMedia(client.wac, fetchUrl)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = client.SendMessage(e.Info.Chat, msg)

	if err != nil {
		log.Println(err)
		return
	}
}

func (client *PrinceClient) handleMessageChat(e *events.Message) {
	content, _ := utils.GetTextContext(e.Message)

	if e.Info.IsFromMe {
		return
	}

	reply, err := utils.GetChatReponse(e.Info.Chat, content)

	if err != nil {
		log.Println(err)
		return
	}

	_, err = client.SendMessage(e.Info.Chat, &waProto.Message{
		Conversation: &reply,
	})

	if err != nil {
		log.Println(err)
		return
	}
}

func (client *PrinceClient) sendRepeatedMessage(msg db.RepeatedMessage) error {
	jid, err := types.ParseJID(msg.JID)
	if err != nil {
		log.Println("Failed to send message:", err)
		return err
	}

	_, err = client.SendCommandMessage(jid, msg.User, utils.CreateTextMessage(msg.Message))
	if err != nil {
		log.Println("Failed to send message:", err)
		return err
	}

	// Update next date
	switch msg.Repeat {
	case "YEARLY":
		msg.NextDate = msg.NextDate.AddDate(1, 0, 0)
	case "MONTHLY":
		msg.NextDate = msg.NextDate.AddDate(0, 1, 0)
	case "WEEKLY":
		msg.NextDate = msg.NextDate.AddDate(0, 0, 7)
	case "DAILY":
		msg.NextDate = msg.NextDate.AddDate(0, 0, 1)
	}

	err = db.UpdateNextDate(msg.ID, msg.NextDate)
	if err != nil {
		return err
	}

	return nil
}

func (client *PrinceClient) sendRepeatedMessages() {
	msgs, err := db.GetRepeatedMessageToday()
	if err != nil {
		log.Println("Failed to send repeated messages:", err)
		return
	}

	for _, msg := range msgs {
		go client.sendRepeatedMessage(msg)
	}

	log.Println("Sent", len(msgs), "repeated messages")
}
