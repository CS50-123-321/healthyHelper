package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	every("*/1 * * * *", func() { bot.Remind("Daily Reminder!!!") })
	//every("0 7 * * *", func() { bot.SetOffDay() }) // Runs daily at 7 AM
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
