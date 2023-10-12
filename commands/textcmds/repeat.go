package textcmds

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"

	"github.com/SushiWaUmai/prince/db"
	"github.com/SushiWaUmai/prince/utils"
)

func RepeatCommand(client *whatsmeow.Client, chat types.JID, user string, ctx *waProto.ContextInfo, pipe *waProto.Message, args []string) (*waProto.Message, error) {
	// 1. arg: start date xx.xx.xxxx
	// 2. arg: repeat "YEARLY","MONTHLY","WEEKLY","DAILY"
	// 3-n. arg: command
	if len(args) < 3 {
		return utils.CreateTextMessage("Usage: repeat <start date> <repeat> <command>"), errors.New("Not enough arguments")
	}

	// Get the date
	var date time.Time
	var err error
	if strings.ToUpper(args[0]) == "TODAY" {
		date = time.Now()
	} else { 
		date, err = time.Parse("02.01.2006", args[0])
		if err != nil {
			return utils.CreateTextMessage("Error parsing date. Please use format dd.mm.yyyy"), err
		}
	}

	// Get the repeat
	repeat := strings.ToUpper(args[1])
	if (repeat != "YEARLY") && (repeat != "MONTHLY") && (repeat != "WEEKLY") && (repeat != "DAILY") {
		return utils.CreateTextMessage("Error parsing repeat. Please use one of 'YEARLY', 'MONTHLY', 'WEEKLY' or 'DAILY'"), errors.New("Could not parse repeat argument")
	}

	// Get the command 
	command := strings.Join(args[2:], " ")

	// Save the command
	_, err = db.CreateRepeatedCommand(chat.String(), user, command, repeat, date)
	if err != nil {
		return nil, err
	}

	return utils.CreateTextMessage(fmt.Sprintf("Saved! Executing first command at %s", date.Format("02.01.2006"))), nil
}

func init() {
	utils.CreateCommand("repeat", "USER", "Creates a repeated command", RepeatCommand)
}
