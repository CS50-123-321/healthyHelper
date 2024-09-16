package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute
	//every("0 7 * * *", func() { bot.Iterator() }) // Runs daily at 7 AM, set false for the day
	every("* * * * *", func() { bot.Iterator() })
	bot.BotInit()
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
