package api

import (
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v3"
)

func InitRoutes() {
	LunchMiniApp()
}
func server(TId int64) {
	var h bot.Habit
	var err error

	// Initialize Gin router
	router := gin.Default()

	router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static/")

	// Serve the index.html file at the root ("/")
	router.GET("", func(c *gin.Context) {
		log.Println("Get / context:", c)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "books",
		})
	})

	router.GET("/create-habit", func(c *gin.Context) {
		log.Println("GET /create-habit context:", c)
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
		h.TeleID = int(TId)
		err = Create(h)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": h})
		config.B.Close()
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

	if err := router.Run(":9000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
func LunchMiniApp() {
	inlineBtn := tele.InlineButton{
		Text:   "Open Mini App!",
		WebApp: &tele.WebApp{URL: "https://familycody.fly.dev"},
	}

	inlineKeys := [][]tele.InlineButton{
		{tele.InlineButton(inlineBtn)},
	}

	// Set up a handler for messages to send the button
	config.B.Handle("/start", func(c tele.Context) error {
		return c.Send("Click the button below:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
	})
	// Handle callback queries when the button is clicked
	config.B.Handle(&inlineBtn, func(c tele.Context) error {
		user := c.Sender() // Get the user who clicked the button
		userID := user.ID
		server(userID)
		config.B = nil
		return nil
	})
	log.Println("bot is running")
	config.B.Start()
}
