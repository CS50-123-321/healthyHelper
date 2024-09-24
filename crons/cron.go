package cron

import (
	"StreakHabitBulder/bot"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	every12("02:30AM", func() { fmt.Println("hey") })    // 8:15 AM Baghdad time
	every12("08:15AM", func() { bot.Act("MentionAll") }) // 8:15 AM Baghdad time
	every12("08:20AM", func() { bot.Act("bestStreak") }) // 8:20 AM Baghdad time
	every12("07:00PM", func() { bot.Act("SendStatus") }) // 7:00 PM Baghdad time
	every12("10:00PM", func() { bot.Act("dailyWatch") }) // 10:00 PM Baghdad time
}

func every12(twelveHourTime string, job func()) {
	t, err := time.Parse("03:04PM", twelveHourTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	baghdadTime := t.Add(3 * time.Hour) 

	duration := fmt.Sprintf("%d %d * * *", baghdadTime.Minute(), baghdadTime.Hour())

	every(duration, job)
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Executing scheduled job...")
		job()
	})
	c.Start()
}
