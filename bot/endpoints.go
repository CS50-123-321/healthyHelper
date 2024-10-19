package bot

import (
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func StartBot() {
	StreakListner()
}

func StreakListner() {
	config.B.Handle(tele.OnVideoNote, func(c tele.Context) (err error) {
		VideoImgListner(c)
		return
	})
	config.B.Handle(tele.OnPhoto, func(c tele.Context) (err error) {
		VideoImgListner(c)
		return
	})

	log.Println("Listeners are running")
	config.B.Start()
}
func VideoImgListner(c tele.Context) (err error) {
	var h Habit

	h.TeleID = int(c.Sender().ID)
	log.Println("logged used", h.TeleID)
	if c.Chat().Type == tele.ChatGroup || c.Chat().Type == tele.ChatSuperGroup { // this is only if the user is adding the mini app to another group
		h.GroupId = int(c.Chat().ID)
	}
	h.TeleID = 6778676335
	log.Println("Edited used", h.TeleID)
	var key = RK(h.GroupId, h.TeleID)
	log.Println("Key used", key)
	h, err = GetMemberHabit(key)
	log.Println("habit used", h)
	if err != nil {
		return fmt.Errorf("error getting days record: %v", err)
	}
	if h.TeleID == 0 {
		msg := fmt.Sprintf("%v hasn't made a habit to commit to within a group", c.Sender().FirstName)
		config.B.Reply(c.Message(), msg)
		return fmt.Errorf("this teleID doesn''t exsist in redis: %v", err)
	}

	err = json.Unmarshal([]byte(h.DaysLogByte), &h.DaysLog)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	h.DaysLog[time.Now().Format("2006-01-02")] = true

	h.DaysLogByte, err = json.Marshal(h.DaysLog)
	if err != nil {
		return err
	}
	h.SetUserStreak()
	if h.TopHit == 0 {
		h.TotalDays = 0
		h.Streaked = 0
		LevelMessage(h, 0)
		return config.Rdb.HSet(context.Background(), key, h).Err()
	}
	log.Println("total and committed used", h.TotalDays, h.CommitmentPeriod)
	// handle when the user finished the period
	if h.TotalDays == h.CommitmentPeriod {
		log.Println("You have made it, set another challege and start again!!")
		config.B.Reply(c.Message(), "You have made it, set another challege and start again!!")
		return nil
	}
	err = config.Rdb.HSet(context.Background(), key, h).Err()
	if err != nil {
		return err
	}

	// setting The Current level and informing them.
	dump := make(map[int]Habit)
	err = json.Unmarshal(h.NotificationLogBytes, &h.NotificationLog)
	if err != nil {
		return err
	}
	dump[h.TeleID] = h
	SetMemberLevel(dump)
	percentageCompleted := (h.TotalDays * 100 / h.CommitmentPeriod)
	level := GetHabitLevel(percentageCompleted)
	// replay with the day count.
	msg := fmt.Sprintf("%v :Day %d, Streak: %d", level, h.TotalDays, h.Streaked)

	var hm HabitMessage
	hm.HabitMsgs(Habit{}, "")
	Aimsg, err := GenerateText(hm.InstantReply.AfterDayCounter)
	if err != nil {
		return err
	}

	config.B.Reply(c.Message(), msg)
	config.B.Reply(c.Message(), EscapeMarkdown(Aimsg))

	return err
}
