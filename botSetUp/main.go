package main

import (
	"StreakHabitBulder/exec"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	log.Println("Server is running!")
	godotenv.Load()
	exec.Init()
}
