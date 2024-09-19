package main

import (
	exec "StreakHabitBulder/exec"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Server is running!")
	godotenv.Load()
	exec.Init()

}
