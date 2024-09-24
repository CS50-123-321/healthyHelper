package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	//every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute for testing.
	every("16 16 * * *", func() { bot.Remind("Let's keep our streak on!!||صباح الخير مو تنسون الرياضة||") }) // 8:15 AM Baghdad time
	every("16 17 * * *", func() { bot.Act("bestStreak") })                                                  // 8:20 AM Baghdad time
	every("0 16 * * *", func() { bot.Act("SendStatus") })                                                   // 7:00 PM Baghdad time
	every("0 18 * * *", func() { bot.Act("dailyWatch") })                                                   // 10:00 PM Baghdad time
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
