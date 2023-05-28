package client

import (
	"log"

	"github.com/SushiWaUmai/prince/db"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func (client *PrinceClient) handleMessage(e *events.Message) {
	client.handleCommand(e.Message, e.Info.ID, e.Info.Chat, e.Info.Sender.User)
}

func (client *PrinceClient) sendRepeatedMessage(msg db.RepeatedMessage) error {
	jid, err := types.ParseJID(msg.JID)
	if err != nil {
		log.Println("Failed to send message:", err)
		return err
	}
	_, err = client.SendCommandMessage(jid, msg.User, &waProto.Message{
		Conversation: proto.String(msg.Message),
	})
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

	db.UpdateNextDate(msg.ID, msg.NextDate)

	return nil
}

func (client *PrinceClient) sendRepeatedMessages() {
	msgs := db.GetRepeatedMessageToday()

	for _, msg := range msgs {
		go client.sendRepeatedMessage(msg)
	}

	log.Println("Sent", len(msgs), "repeated messages")
}
