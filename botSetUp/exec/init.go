package exec

import (
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	cron "StreakHabitBulder/crons"
)

func Init() {
	config.InitTele()
	config.InitRedis()
	bot.Iterator() // for testingjn
	bot.BotInit()
	cron.InitCron()
}
