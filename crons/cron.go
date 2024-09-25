package cron

import (
	"StreakHabitBulder/bot"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)


func InitCron() {
	bot.Act("GenerateAiRandomMember")
	//msg, _ := plugin.GeneratePersonlizeResponse(habit)
	every12("03:51AM", func() { fmt.Println("hey") })                // 08:15 AM Baghdad Time
	every12("02:50PM", func() { bot.Act("GenerateAiRandomMember") }) // 08:15 AM Baghdad Time
	every12("08:10AM", func() { bot.Act("MentionAll") })             // 08:15 AM Baghdad Time
	every12("08:20AM", func() { bot.Act("bestStreak") })             // 08:20 AM Baghdad Time
	every12("07:00PM", func() { bot.Act("SendStatus") })             // 07:00 PM Baghdad Time
	every12("10:00PM", func() { bot.Act("dailyWatch") })             // 10:00 PM Baghdad Time
}

func every12(twelveHourTime string, job func()) {
	t, err := time.Parse("03:04PM", twelveHourTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	baghdadToUTC := t.Add(-3 * time.Hour)

	duration := fmt.Sprintf("%d %d * * *", baghdadToUTC.Minute(), baghdadToUTC.Hour())

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
