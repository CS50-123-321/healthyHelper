package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	//every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute
	every("0 7 * * *", func() { bot.Act("SetDayOff") }) // Runs daily at 7 AM, sends Analytics Message.
	every("2 0 * * *", func() { bot.Act("SendStatus") })
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
