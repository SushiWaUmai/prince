package env

import (
	"os"
)

var (
	BOT_PREFIX     string
	OPENAI_API_KEY string
)

func loadEnv() {
	BOT_PREFIX = os.Getenv("BOT_PREFIX")
	OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
}
