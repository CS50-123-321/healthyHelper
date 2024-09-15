package api

import (
	"familyFormUi/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tele "gopkg.in/telebot.v3"
)

func InitRoutes() {
	GetUserBotID()
}
func server(TId int64) {
	var h *Habit

	// Initialize Gin router
	router := gin.Default()

	router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static/")

	// Serve the index.html file at the root ("/")
	router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"books": "books",
		})
	})

	// Define the route to handle habit form submissions
	router.POST("/save-habit", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&h); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Save the h data in Redis
		h.TeleID = TId
		err := h.Create()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": h})
	})

	router.Run(":9000")
}
func GetUserBotID() {
	config.B.Handle("/start", func(c tele.Context) error {
		c.Send("Click Join!!")
		server(c.Sender().ID)
		return nil
	})
	log.Println("bot is running")
	config.B.Start()
}
