package be

import (
	"StreakHabitBulder/config"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Habit struct {
	Name             string `json:"name" redis:"name" binding:"required"`
	HabitName        string `json:"habitName" redis:"habit_name" binding:"required"`
	CommitmentPeriod string `json:"commitmentPeriod"  redis:"commitment_Period" binding:"required"`
}

func InitRoutes() {
	// Initialize Gin router
	router := gin.Default()
	router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static")
	router.GET("/t", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "books",
		})
	})

	// Define the route to handle habit form submissions
	router.POST("/save-habit", func(c *gin.Context) {
		var h Habit
		if err := c.ShouldBindJSON(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Save the h data in Redis
		err := h.Add()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Habit saved successfully!"})
	})

	// Start the server
	log.Println("Server starting on :8080")
	router.Run(":8080")
}

func (h *Habit) Add() error {
	key := fmt.Sprintf("habitMember:%d", 12)

	err := config.Rdb.HSet(context.Background(), key, h).Err()
	if err != nil {
		return err
	}
	return err
}
