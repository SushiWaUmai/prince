package env

import (
	"os"
)

var (
	BOT_PREFIX                byte
	OPENAI_API_KEY            string
	VOICEVOX_ENDPOINT         string
	STABLE_DIFFUSION_ENDPOINT string
)

func loadEnv() {
	if os.Getenv("BOT_PREFIX") != "" {
		BOT_PREFIX = os.Getenv("BOT_PREFIX")[0]
	}

	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	VOICEVOX_ENDPOINT = os.Getenv("VOICEVOX_ENDPOINT")
	STABLE_DIFFUSION_ENDPOINT = os.Getenv("STABLE_DIFFUSION_ENDPOINT")
}
