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
		// err := c.ShouldBindQuery(&p)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
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

	// Set up a handler for messages to send the button
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
		//check if the user has a habit made before or not
		// h, err := GetHabit(int(c.Sender().ID))
		// if err != nil {
		// 	return err
		// }
		// if h.TeleID != 0 {
		// 	c.Send(fmt.Sprintf("You already have a habit: %v", h))
		// }

		// Send the habit information separately
		c.Send("Click the button below:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
		//config.B = nil
		log.Println("Stopping the bot")
		return nil
	})
	// Handle callback queries when the button is clicked
	// config.B.Handle(&inlineBtn, func(c tele.Context) error {
	// 	user := c.Sender() // Get the user who clicked the button
	// 	userID := user.ID
	// 	log.Println("Stopping the bot")
	// 	return nil
	// })
	config.B.Start()
}

func GetHabit(id int) (h bot.Habit, err error) {
	return bot.GetDaysRecord(bot.RK(id))
}
