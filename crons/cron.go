package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	//every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute
	//every("15 5 * * *", func() { bot.Act("SetDayOff") })                  // 8:15 AM Baghdad time
	every("15 5 * * *", func() { bot.Remind("مو تنسون تسوون رياضة اليوم") }) // 8:15 AM Baghdad time
	every("0 16 * * *", func() { bot.Act("SendStatus") })                    // 7:00 PM Baghdad time
	every("0 18 * * *", func() { bot.Act("dailyWatch") })                    // 10:00 PM Baghdad time
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}

