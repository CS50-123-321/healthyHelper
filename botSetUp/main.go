package main

import (
	exec "StreakHabitBulder/exec"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	exec.Init()

}
