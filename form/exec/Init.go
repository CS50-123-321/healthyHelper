package exec

import (
	"familyFormUi/api"
	"familyFormUi/bot"
	"familyFormUi/config"
)

func Init() {
	config.InitRedis()
	config.InitTele()
	api.InitRoutes()
	bot.StreakListner()
}
