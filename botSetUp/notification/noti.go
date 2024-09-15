package notif

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitNotification() {
	every("*/1 * * * *", func() { bot.Remind("Daily Reminder!!!") })
	every("0 7 * * *", func() { bot.InitOffDay() }) // Runs daily at 7 AM
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}

// 	c := cron.New()
// 	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
// 	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
// 	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
// 	c.Start()
// }
