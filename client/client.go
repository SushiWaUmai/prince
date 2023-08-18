package client

import (
	"context"
	"log"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/robfig/cron/v3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
)

type PrinceClient struct {
	wac            *whatsmeow.Client
	cronJob        *cron.Cron
	eventHandlerId uint32
	commandPrefix  byte
}

func CreatePrinceClient(prefix byte, deviceStore *store.Device) *PrinceClient {
	clientLog := waLog.Stdout("Client", "INFO", true)
	wac := whatsmeow.NewClient(deviceStore, clientLog)
	client := &PrinceClient{
		wac:            wac,
		eventHandlerId: 0,
		commandPrefix:  prefix,
		cronJob:        cron.New(),
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

	err := client.startCronJobs()
	if err != nil {
		log.Println(err)
	}
}

func (client *PrinceClient) Disconnect() {
	log.Println("Stopping Cron Job...")
	client.cronJob.Stop()
	client.wac.Disconnect()
}

func (client *PrinceClient) SendMessage(chat types.JID, msg *waProto.Message) (resp whatsmeow.SendResponse, err error) {
	return client.wac.SendMessage(context.Background(), chat, msg)
}

func (client *PrinceClient) SendCommandMessage(chat types.JID, user string, msg *waProto.Message) (resp whatsmeow.SendResponse, err error) {
	resp, err = client.SendMessage(chat, msg)
	if err != nil {
		return resp, err
	}

	client.handleCommand(msg, resp.ID, chat, user)
	return resp, err
}

func (client *PrinceClient) register() {
	client.eventHandlerId = client.wac.AddEventHandler(client.eventHandler)
}

func (client *PrinceClient) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go client.handleMessage(v)
	}
}

func (client *PrinceClient) startCronJobs() error {
	_, err := client.cronJob.AddFunc("0 0 * * *", client.sendRepeatedMessages)

	if err != nil {
		return err
	}

	log.Println("Starting Cron Job...")
	client.cronJob.Start()

	return nil
}
