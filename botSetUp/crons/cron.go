package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	//every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute
	every("0 7 * * *", func() { bot.Act("SetDayOff") })   // Runs daily at 7 AM, sends Analytics Message.
	every("0 19 * * *", func() { bot.Act("SendStatus") }) // Send at 7 PM daily
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
