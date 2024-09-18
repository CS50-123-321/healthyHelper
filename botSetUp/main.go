package main

import (
	exec "StreakHabitBulder/exec"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	exec.Init()

}
