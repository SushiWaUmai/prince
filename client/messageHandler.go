package client

import (
	"log"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/env"
	"github.com/SushiWaUmai/prince/utils"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (client *PrinceClient) handleMessage(e *events.Message) {
	client.handleCommand(e.Message, e.Info.ID, e.Info.Chat, e.Info.Sender.User)
	client.handleMessageEvents(e)
}

func (client *PrinceClient) handleMessageEvents(e *events.Message) {
	jid := e.Info.Chat.String()

	msgEvents := db.GetMessageEvents(jid)

	ctx := &waProto.ContextInfo{
		QuotedMessage: e.Message,
	}

	chat := e.Info.Chat
	user := client.wac.Store.ID.User

	for _, evt := range msgEvents {
		result, err := RunCommand(client.wac, string(env.BOT_PREFIX)+evt.Content, ctx, chat, user)

		if result != nil && err == nil {
			client.SendCommandMessage(chat, user, result)
		}
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
	msgs := db.GetRepeatedMessageToday()

	for _, msg := range msgs {
		go client.sendRepeatedMessage(msg)
	}

	log.Println("Sent", len(msgs), "repeated messages")
}
