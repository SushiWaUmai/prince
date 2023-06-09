package textcmds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func init() {
	utils.CreateCommand("repeat", "USER", func(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
		// 1. arg: start date xx.xx.xxxx
		// 2. arg: repeat "YEARLY","MONTHLY","WEEKLY","DAILY"
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
		repeat := strings.ToUpper(args[1])
		if (repeat != "YEARLY") && (repeat != "MONTHLY") && (repeat != "WEEKLY") && (repeat != "DAILY") {
			response := &waProto.Message{
				Conversation: proto.String("Error parsing repeat. Please use one of 'YEARLY', 'MONTHLY', 'WEEKLY' or 'DAILY'"),
			}
			return response, errors.New("Could not parse repeat")
		}

		// Get the message
		message := strings.Join(args[2:], " ")

		// Save the message
		db.CreateRepeatedMessage(chat.String(), user, message, repeat, date)

		// Reply
		response := &waProto.Message{
			Conversation: proto.String(fmt.Sprintf("Saved! Sending first message at %s", date.Format("02.01.2006"))),
		}

		return response, nil
	})
}
