package api

import (
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
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
		if h.SoloSre != "" && h.SoloSre == "no" {
			h.IsGroup = false
		} else if h.SoloSre != "" && h.SoloSre == "yes" {
			h.IsGroup = true
		}
		if h.IsGroup {
			// fetch the groups ids
			// Save maps group id with its members.
			chatMember, err := config.B.ChatMemberOf(tele.ChatID(-1002239647108), tele.ChatID(175864127))
			log.Println("------------", chatMember.User, err)
			msg := "🚀 <a href='https://t.me/StreakForBetterHabits_Bot?startgroup=true'>Click here to add the bot to your group</a> and let it track everyone's progress!"
			_, err = config.B.Send(tele.ChatID(h.TeleID), msg, tele.ModeHTML)
			if err != nil {
				log.Println("err in /create-habit", err)
				return
			}
		}
		return
		err = Create(h)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": h})
	})

	router.GET("/progress", func(c *gin.Context) {
		var p ProgresRequest
		tid := c.Query("tele_id")
		p.TeleID, err = strconv.Atoi(tid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
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
	router.GET("/dashboard", func(c *gin.Context) {
		AllMemberHabits := bot.Act(" ")
		// sortedMembers :=
		sort.Slice(AllMemberHabits, func(i, j int) bool {
			return AllMemberHabits[i].TotalDays > AllMemberHabits[j].TotalDays
		})
		c.HTML(http.StatusOK, "admin.html", gin.H{
			"Habit": AllMemberHabits,
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
		log.Println("start command is running")
		if c.Chat().Type == tele.ChatGroup || c.Chat().Type == tele.ChatSuperGroup { // this is only if the user is adding the mini app to another group
			groupID := c.Chat().ID
			userID := c.Sender().ID
			log.Println("saving to redis")
			return SaveGroupIDToRedis(int(userID), int(groupID))
		}
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
		webAppURL := fmt.Sprintf("https://familycody.fly.dev/progress?tele_id=%d", c.Sender().ID)
		inlineBtn := tele.InlineButton{
			Text:   "Your Progress!",
			WebApp: &tele.WebApp{URL: webAppURL},
		}
		inlineKeys := [][]tele.InlineButton{
			{inlineBtn},
		}
		c.Send("Click the button below:", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
		return nil
	})
	config.B.Handle("/Admin", func(c tele.Context) error {
		log.Println("bot is running")
		webAppURL := "https://familycody.fly.dev/dashboard"
		inlineBtn := tele.InlineButton{
			Text:   "Dashboard!",
			WebApp: &tele.WebApp{URL: webAppURL},
		}
		inlineKeys := [][]tele.InlineButton{
			{inlineBtn},
		}
		c.Send("Admin Dashboard", &tele.ReplyMarkup{InlineKeyboard: inlineKeys})
		return nil
	})
	config.B.Start()
}

func GetHabit(id int) (h bot.Habit, err error) { return bot.GetDaysRecord(bot.RK(id)) }
func SaveGroupIDToRedis(userid, groupId int) (err error) {
	err = config.Rdb.HSet(context.Background(), fmt.Sprintf("habitMember:%v", userid), "groupID", groupId).Err()
	if err != nil {
		return
	}
	err = config.Rdb.SAdd(context.Background(), "groupIds", groupId).Err()
	if err != nil {
		return
	}
	err = config.Rdb.SAdd(context.Background(), fmt.Sprintf("habitByGroup:%v", groupId), userid).Err()
	return
}
