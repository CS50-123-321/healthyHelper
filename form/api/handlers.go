package api

import (
	"context"
	"encoding/json"
	"familyFormUi/config"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

func (h *Habit) Create() (err error) {
	switch h.CommitmentPeriodStr {
	case "10 days":
		h.CommitmentPeriod = 10
	case "20 days":
		h.CommitmentPeriod = 20
	case "30 days":
		h.CommitmentPeriod = 30
	}

	h.DaysLog = make(map[int]bool)
	h.NotificationLog = make(map[int]bool)
	h.TotalDays = 0
	h.TopHit = 0
	h.Streaked = 0
	daysLogJSON, err := json.Marshal(h.DaysLog)
	if err != nil {
		return err
	}
	h.NotificationLogBytes, err = json.Marshal(h.DaysLog)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("habitMember:%d", h.TeleID)
	err = config.Rdb.HSet(context.Background(), key, map[string]interface{}{
		"name":              h.Name,
		"habit_name":        h.HabitName,
		"commitment_period": h.CommitmentPeriod,
		"tele_id":           h.TeleID,
		"streaked":          h.Streaked,
		"days_log":          daysLogJSON, // Store as JSON string
		"total_days":        h.TotalDays,
		"top_hit":           h.TopHit,
		"notification_log":  h.NotificationLogBytes,
	}).Err()
	if err != nil {
		return err
	}
	// Adding all client ID to one space.
	err = config.Rdb.ZAdd(context.Background(), "MembersIDS", redis.Z{
		Score:  float64(time.Now().Unix()), // TODO: Make it tele group id tocatogrize them.
		Member: h.TeleID,
	}).Err()
	if err != nil {
		return err
	}

	// //TODO: notify
	// msg := fmt.Sprintf("%v, %v", h, h.TeleID)
	// botID, err := strconv.Atoi(os.Getenv("TestingBotID"))
	// if err != nil {
	// 	log.Println("err:", err)
	// 	return
	// }
	// t, err := config.B.Send(tele.ChatID(botID), msg)
	// if err != nil {
	// 	log.Println("err:", err)
	// 	return
	// }
	// log.Println(t)
	return
}
