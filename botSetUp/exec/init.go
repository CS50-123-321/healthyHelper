package exec

import (
	"StreakHabitBulder/config"
	cron "StreakHabitBulder/crons"
)

func Init() {
	config.InitTele()
	config.InitRedis()
	// go func() {
	// 	bot.BotInit()
	// 	select {}
	// }()
	go func() {
		cron.InitCron()
		select {}
	}()
	select {}
}
