package rest

import (
	"StreakHabitBulder/DB"
	"StreakHabitBulder/bot"
	notif "StreakHabitBulder/notification"
)

func Init() {
	bot.InitTele()
	DB.InitRedis()
	notif.InitNotification()
	bot.InitOffDay()
}
