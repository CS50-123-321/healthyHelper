package cron

import (
	"StreakHabitBulder/bot"
	"fmt"

	"github.com/robfig/cron/v3"
)

func InitCron() {
	//every("*/1 * * * *", func() { bot.Remind("Minute Reminder!!!") }) // Runs every minute
	every("15 8 * * *", func() { bot.Act("SetDayOff") })                    // Runs daily at 7 AM, sends Analytics Message.
	every("15 8 * * *", func() { bot.Remind("مو تنسون تسوون رياضة اليوم") }) // Runs daily at 7 AM, sends Analytics Message.
	every("0 19 * * *", func() { bot.Act("SendStatus") })                   // Send at 7 PM daily
	every("0 22 * * *", func() { bot.Act("ShitOn") })                       // Send at 1- PM daily, easu, shit on the members who don't stick to the habit!! "0 22 * * *"
}

func every(duration string, job func()) {
	c := cron.New()
	c.AddFunc(duration, func() {
		fmt.Println("Every 6 hour thirty")
		job()
	})
	c.Start()
}
// FlyV1 fm2_lJPECAAAAAAAB6SsxBB3zw6+EK2v5iRYwlUkut7lwrVodHRwczovL2FwaS5mbHkuaW8vdjGWAJLOAAyDTh8Lk7lodHRwczovL2FwaS5mbHkuaW8vYWFhL3YxxDx3DLNaBKaCfIUOKt1FJYQVsCCkHhyS4gL6mt64jP5ojrhbZvOaMtqoHN4GC5COIrLlyvBJ9k2GKiky/f/ETsZWTXDZx6czd72NCqqgEjDdUXUrFybtQX+809I/biCOb6VlLNrNQP90Sf90diVHVgsJysZzMfLKWnIgZHkHdMH89Q6PgP/coGlKRRnN4w2SlAORgc4ARY3uHwWRgqdidWlsZGVyH6J3Zx8BxCAVHi9F7oQXIL+J8hdm0rdvOFjWijIMc8gLZNTU4KudDg==,fm2_lJPETsZWTXDZx6czd72NCqqgEjDdUXUrFybtQX+809I/biCOb6VlLNrNQP90Sf90diVHVgsJysZzMfLKWnIgZHkHdMH89Q6PgP/coGlKRRnN48QQhv7WbTqRmLylUPu8Hn/es8O5aHR0cHM6Ly9hcGkuZmx5LmlvL2FhYS92MZgEks5m76dizwAAAAE9gz1wF84ADCPsCpHOAAwj7AzEELeU58XBKPH5ntEdwn8aAnLEIBCOTMJetPdEZRKGdmOw5cYkgjgO+rLC9bi2d0nZ4qxZ