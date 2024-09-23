package exec

import (
	"StreakHabitBulder/api"
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	cron "StreakHabitBulder/crons"
	"fmt"
	"sync"
)

func Init() {
	// Control CronJobs
	config.InitTele()
	config.InitRedis()
	waitGroup()
}
func waitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Running Cron!")
		cron.InitCron()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Running BotInit!")
		bot.StartBot()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("running server")
		api.Server()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Lunching mini app")
		api.LunchMiniApp()
	}()
	wg.Wait()
}
