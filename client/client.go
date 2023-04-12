package client

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/SushiWaUmai/prince/commands"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

type PrinceClient struct {
	wac            *whatsmeow.Client
	eventHandlerId uint32
	commandPrefix  string
}

func CreatePrinceClient(prefix string, deviceStore *store.Device) *PrinceClient {
	clientLog := waLog.Stdout("Client", "ERROR", true)
	wac := whatsmeow.NewClient(deviceStore, clientLog)
	client := &PrinceClient{
		wac:            wac,
		eventHandlerId: 0,
		commandPrefix:  prefix,
	}

	client.register()
	return client
}

func (client *PrinceClient) Connect() {
	if client.wac.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.wac.GetQRChannel(context.Background())
		err := client.wac.Connect()
		if err != nil {
			log.Fatalln(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				log.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err := client.wac.Connect()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (client *PrinceClient) Disconnect() {
	client.wac.Disconnect()
}

func (client *PrinceClient) register() {
	client.eventHandlerId = client.wac.AddEventHandler(client.eventHandler)
}

func (client *PrinceClient) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		client.handleCommand(v)
	}
}

func (client *PrinceClient) handleCommand(e *events.Message) {
	msg := e.Message

	if !e.Info.IsFromMe {
		return
	}

	content := msg.GetConversation()

	if !strings.HasPrefix(content, client.commandPrefix) {
		return
	}

	content = content[len(client.commandPrefix):]

	// split the command name with the arguments
	split := strings.SplitN(content, " ", -1)
	cmdName := split[0]
	cmdArgs := split[1:]

	log.Println("Runnning commmand", cmdName, "with args", cmdArgs)
	for _, cmd := range commands.CommandList {
		if cmdName == cmd.Name {
			cmd.Execute(client.wac, e, cmdArgs)
		}
	}
}
