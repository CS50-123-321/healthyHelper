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
	var key = RK(h.TeleID)
	// Get old records
	h, err = GetDaysRecord(key)
	if err != nil {
		return fmt.Errorf("error getting days record: %v", err)
	}
	if h.TeleID == 0 {
		msg := fmt.Sprintf("%v hasn't made a habit to commit to", c.Sender().FirstName)
		config.B.Reply(c.Message(), msg)
		return fmt.Errorf("this teleID doesn''t exsist in redis: %v", err)
	}

	err = json.Unmarshal([]byte(h.DaysLogByte), &h.DaysLog)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	h.DaysLog[time.Now().Day()] = true

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
	config.B.Reply(c.Message(), msg)
	return err
}
