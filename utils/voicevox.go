package utils

import (
	"log"
	"net/url"

	"github.com/SushiWaUmai/prince/env"
	"github.com/aethiopicuschan/voicevox"
)

func CreateVoiceVoxClient() *voicevox.Client {
	url, err := url.Parse(env.VOICEVOX_ENDPOINT)
	if err != nil {
		log.Println("Failed to register zundamon command:", err)
		return nil
	}

	voicevox := voicevox.NewClient(url.Scheme, url.Host)
	return voicevox
}
