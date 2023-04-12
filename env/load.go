package env

import (
	"os"
)

var (
	BOT_PREFIX = "!"
)

func loadEnv() {
	BOT_PREFIX = os.Getenv("BOT_PREFIX")
}
