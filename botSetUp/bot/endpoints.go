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
	var h Habit
	config.B.Handle(tele.OnText, func(c tele.Context) (err error) {
		h.TeleID = int(c.Sender().ID)
		var key = RK(h.TeleID)
		// Get old record
		h, err = getDaysRecord(key)
		if err != nil {
			return fmt.Errorf("error getting days record: %v", err)
		}
		if h.TeleID == 0 {
			msg := fmt.Sprintf("%v hasn't made a habit to commit to", c.Sender().FirstName)
			config.B.Reply(c.Message(), msg)
			return fmt.Errorf("this teleID doesn''t exsist in redis: %v", err)
		}

		// Unmarshell it to the struct
		err = json.Unmarshal([]byte(h.DaysLogByte), &h.DaysLog)
		if err != nil {
			return fmt.Errorf("error unmarshalling JSON: %v", err)
		}
		// Marking day as true
		h.DaysLog[time.Now().Day()] = true

		h.DaysLogByte, err = json.Marshal(h.DaysLog)
		if err != nil {
			return err
		}
		h.SetUserStreak() //calc streak info
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
		// replay with the day count.
		msg := fmt.Sprintf("Day %d, Streak: %d", h.TotalDays, h.Streaked)
		config.B.Reply(c.Message(), msg)
		// setting The Current level and informing them.
		dump := make(map[int]Habit)
		err = json.Unmarshal(h.NotificationLogBytes, &h.NotificationLog)
		if err != nil {
			return err
		}
		dump[h.TeleID] = h
		SetMemberLevel(dump)
		return err
	})
	// config.B.Handle("/t", func(c tele.Context) (err error) {
	// 	return c.Send("ds")
	// })
	log.Println("Listeners are running")
	config.B.Start()
}
