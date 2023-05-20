package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
)

func init() {
	createCommand("repeat", func(client *whatsmeow.Client, messageEvent *events.Message, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		// 1. arg: start date xx.xx.xxxx
		// 2. arg: repeat "Yearly","Monthly","Weekly","Daily"
		// 3-n. arg: message
    // TODO: use pipe
		if len(args) < 3 {
			repsonse := &waProto.Message{
				Conversation: proto.String("Usage: repeat <start date> <repeat> <message>"),
			}
			return repsonse, errors.New("Not enough arguments")
		}

		// Get the date
		date, err := time.Parse("02.01.2006", args[0])
		if err != nil {
			response := &waProto.Message{
				Conversation: proto.String("Error parsing date. Please use format dd.mm.yyyy"),
			}
			return response, err
		}

		// Get the repeat
		repeat := args[1]
		if (repeat != "y") && (repeat != "m") && (repeat != "w") && (repeat != "d") {
			response := &waProto.Message{
				Conversation: proto.String("Error parsing repeat. Please use one of 'y', 'm', 'w' or 'd'"),
			}
			return response, errors.New("Could not parse repeat")
		}

		// Get the message
		message := strings.Join(args[2:], " ")

		// Save the message
		db.CreateRepeatedMessage(messageEvent.Info.Chat.String(), message, repeat, date)

		// Reply
		response := &waProto.Message{
			Conversation: proto.String(fmt.Sprintf("Saved! Sending first message at %s", date.Format("02.01.2006"))),
		}

		return response, nil
	})
}
