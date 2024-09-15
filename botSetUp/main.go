package main

import (
	"StreakHabitBulder/bot"
	rest "StreakHabitBulder/rest"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rest.Init()
	bot.Test()
}
