package bot

import (
	"context"
	"encoding/json"
	"familyFormUi/api"
	"familyFormUi/config"
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

// TODO: What is multiple times per day?
// Make an account
func StreakListner() {
	var h *api.Habit
	config.B.Handle(tele.OnVideoNote, func(c tele.Context) error {
		h.TeleID = c.Sender().ID
		key := fmt.Sprintf("habitMember:%d", h.TeleID)
		h.DaysLog[int(time.Now().Day())] = true
		daysLogJson, err := json.Marshal(h.DaysLog)
		if err != nil {
			return err
		}
		err = config.Rdb.HSet(context.Background(), key, "days_log", daysLogJson).Err()
		return err
	})
	log.Println("Listeners are running")

	config.B.Start()
}
