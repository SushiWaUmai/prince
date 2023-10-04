package commands

import "github.com/SushiWaUmai/prince/db"

func init() {
	db.CreateAlias("autodownload", "onmessage download")
	db.CreateAlias("autopilot", "onmessage chat")
	db.CreateAlias("zundamon", "voicevox 1")
}
