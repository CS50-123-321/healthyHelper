package bot

import (
	"StreakHabitBulder/DB"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func CreateStreak(memberID int, habit string, days int) *Streak {
	return &Streak{
		MemberID:   memberID,
		Habit:      habit,
		StartDate:  time.Now(),
		Days:       make([]int, days),
		TotalDays:  days,
		CurrentDay: 0,
	}
}

func Test() {
	StreakListner()
}
func StreakListner() {
	var h Habit
	b.Handle(tele.OnText, func(c tele.Context) (err error) {
		h.TeleID = int(c.Sender().ID)
		var key = fmt.Sprintf("habitMember:%d", h.TeleID)
		// Get old record
		h, err = getDaysRecord(key)
		if err != nil {
			return fmt.Errorf("error getting days record: %v", err)
		}

		// Unmarshell it to the struct
		err = json.Unmarshal([]byte(h.DaysLogByte), &h.DaysLog)
		if err != nil {
			return fmt.Errorf("error unmarshalling JSON: %v", err)
		}

		//add the new record
		h.DaysLog[time.Now().Day()] = true

		if h.TopHit == 0 {
			h.TotalDays = 0
			h.Streaked = 0
			return DB.Rdb.HSet(context.Background(), key, h).Err()
		}
		h.getStreakByUser()
		h.DaysLogByte, err = json.Marshal(h.DaysLog)
		if err != nil {
			return err
		}
		return DB.Rdb.HSet(context.Background(), key, h).Err()
	})

	log.Println("Listeners are running")
	b.Start()
}
