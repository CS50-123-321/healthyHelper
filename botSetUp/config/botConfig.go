package config

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

var B *tele.Bot

func InitTele() {

	pref := tele.Settings{
		Token:  os.Getenv("TELE_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	var err error
	B, err = tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

}
