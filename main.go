package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow/store/sqlstore"

	"github.com/SushiWaUmai/prince/api"
	"github.com/SushiWaUmai/prince/client"
	_ "github.com/SushiWaUmai/prince/commands"
	"github.com/SushiWaUmai/prince/env"
	_ "github.com/mattn/go-sqlite3"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	api.CreateAPI()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:./data/auth.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to Database")

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalln(err)
	}

	princeClient := client.CreatePrinceClient(env.BOT_PREFIX, deviceStore)
	log.Println("Initialized new Whatsapp Client")

	princeClient.Connect()

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	princeClient.Disconnect()
}
