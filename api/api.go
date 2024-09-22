package api

import (
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v3"
)

func Server() {
	var h bot.Habit
	var err error
	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static/")

	// Serve the index.html file at the root ("/")
	router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "books",
		})
	})

	router.GET("/create-habit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "books",
		})
	})
	// Define the route to handle habit form submissions
	router.POST("/create-habit", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Save the h data in Redis
		h.TeleID, err = strconv.Atoi(h.TeleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session ID"})
			return
		}
		err = Create(h)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": h})
		log.Println("stopping the bot")
	})

	router.GET("/progress", func(c *gin.Context) {
		var p ProgresRequest
		err := c.ShouldBindJSON(&p)
		if err != nil {
			log.Println(err)
			return
		}
		//p.TeleID = 175864127

		err, h := getUserProgress(p.TeleID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user progress"})
			return
		}
		// Pass 'h' (habit data) to the template
		c.HTML(http.StatusOK, "progress.html", gin.H{
			"Habit": h,
		})

	})
	log.Println("Listening on port 0.0.0.0:8888")
	if err := router.Run("0.0.0.0:8888"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func LunchMiniApp() {
	log.Println("bot is running")

	config.B.Handle("/start", func(c tele.Context) error {
		log.Println("bot is running")
		// Create the button with the session ID as a URL parameter
		webAppURL := fmt.Sprintf("https://familycody.fly.dev/create-habit?session=%d", c.Sender().ID)
		inlineBtn := tele.InlineButton{
			Text:   "Open Mini App!",
			WebApp: &tele.WebApp{URL: webAppURL},
		}

		inlineKeys := [][]tele.InlineButton{
			{inlineBtn},
		}
		// Send the habit information separately
		c.Send("Click the button below:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
		log.Println("Stopping the bot")
		return nil
	})
	config.B.Handle("/Me", func(c tele.Context) error {
		log.Println("bot is running")
		// Create the button with the session ID as a URL parameter
		webAppURL := fmt.Sprintf("https://familycody.fly.dev/progress?tele_id=%d", c.Sender().ID)
		inlineBtn := tele.InlineButton{
			Text:   "Your Progress!",
			WebApp: &tele.WebApp{URL: webAppURL},
		}

		inlineKeys := [][]tele.InlineButton{
			{inlineBtn},
		}
		// Send the habit information separately
		c.Send("Click the button below:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
		return nil
	})

	config.B.Start()
}

func GetHabit(id int) (h bot.Habit, err error) {
	return bot.GetDaysRecord(bot.RK(id))
}
