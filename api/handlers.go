package api

import (
	"StreakHabitBulder/bot"
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Create(h bot.Habit) (err error) {

	h.CommitmentPeriod, err = strconv.Atoi(h.CommitmentPeriodStr)
	if err != nil {
		return
	}

	h.DaysLog = make(map[string]bool)
	for i := 0; i < h.CommitmentPeriod; i++ {
		now := time.Now()
		//curMoth := time.Now().Month()
		day := now.AddDate(0, 0, i).Format("2006-01-02")
		h.DaysLog[day] = false
	}
	h.NotificationLog = h.DaysLog // init for daysLog and NotificationLog
	h.TotalDays = 0
	h.TopHit = 0
	h.Streaked = 0
	h.CreatedAt = time.Now()
	h.DaysLogByte, err = json.Marshal(h.DaysLog)
	if err != nil {
		return err
	}
	h.NotificationLogBytes, err = json.Marshal(h.NotificationLog)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("habitMember:%v:%v", h.GroupId, h.TeleID)
	err = config.Rdb.HSet(context.Background(), key, map[string]interface{}{
		"name":              h.Name,
		"habit_name":        h.HabitName,
		"commitment_period": h.CommitmentPeriod,
		"tele_id":           h.TeleID,
		"streaked":          h.Streaked,
		"days_log":          h.DaysLogByte, // Store as JSON string
		"total_days":        h.TotalDays,
		"top_hit":           h.TopHit,
		"notification_log":  h.NotificationLogBytes,
		"created_at":        h.CreatedAt,
		"group_id":          h.GroupId,
	}).Err()
	if err != nil {
		return err
	}
	fmt.Println("adding members")
	// Adding all client ID to one space.
	err = config.Rdb.ZAdd(context.Background(), "MembersIDS", redis.Z{
		Score:  float64(h.GroupId),
		Member: h.TeleID,
	}).Err()
	if err != nil {
		return err
	}
	//msg := fmt.Sprintf("Yay🎉🎉, new habit maker joiner, welcome 🎉%v🎉, Habit: %v", h.Name, h.HabitName)
	//bot.Remind(msg)
	return
}
func getUserProgress(teleId, groupid int) (err error, h bot.Habit) {
	err = config.Rdb.HGetAll(context.Background(), bot.RK(groupid, teleId)).Scan(&h)
	return
}

func SaveGroupIDToRedis(userid, groupId int) (err error) {
	key := fmt.Sprintf("habitMember:%v:%v", groupId, userid)
	err = config.Rdb.HSet(context.Background(), key, "group_id", groupId).Err()
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
