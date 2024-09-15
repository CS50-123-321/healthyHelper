package bot

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	tele "gopkg.in/telebot.v3"
)

var b *tele.Bot

func InitTele() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	pref := tele.Settings{
		Token:  os.Getenv("TELE_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

}
