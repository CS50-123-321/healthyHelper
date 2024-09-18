package bot

import (
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func (m *Members) Add(key string) (err error) {
	err = config.Rdb.HSet(context.Background(), key, m).Err()
	if err != nil {
		return err
	}
	return
}

func getDaysRecord(key string) (h Habit, err error) {
	config.Rdb.HGetAll(context.Background(), key).Scan(&h)
	return h, err
}

func (h *Habit) SetUserStreak() {

	h.TotalDays = 0
	h.Streaked = 0
	// Sorting the days of the daysRecord
	daysLogSlice := make([]int, 0, len(h.DaysLog))
	for day := range h.DaysLog {
		daysLogSlice = append(daysLogSlice, day)

	}
	sort.Ints(daysLogSlice)
	sortedMap := make(map[int]bool)
	for _, day := range daysLogSlice {
		done := h.DaysLog[day]
		if done {
			h.TotalDays++
			h.Streaked++
			if h.TopHit < h.Streaked {
				h.TopHit = h.Streaked
			}
		} else {
			if h.TopHit < h.Streaked {
				h.TopHit = h.Streaked
			}
			h.Streaked = 0
		}
		sortedMap[day] = done
	}
	h.DaysLog = sortedMap
}
func getMembersIDs() (ids []int, err error) {
	err = config.Rdb.ZRange(context.Background(), "MembersIDS", 0, -1).ScanSlice(&ids)
	if err != nil {
		return nil, err
	}
	return
}

func SetMemberLevel(memberHabit map[int]Habit) {
	for _, h := range memberHabit {
		percentageCompleted := 0
		if h.TotalDays > 0 {
			percentageCompleted = (h.TotalDays * 100 / h.CommitmentPeriod)
		}
		ok := h.NotificationLog[time.Now().Day()]
		if !ok { // Send notification only if the user hasn't recevied a notification on this day.
			LevelMessage(h, percentageCompleted)
		} else {
			log.Println("this user has been informed today already")
		}
	}
}

// This get called by the cron job to run daily and sets the day as false, it will be true if the member did sport.

func SetNotificationLog(key string) error {
	h, err := getDaysRecord(key)
	if err != nil {
		log.Println("error getting days record: %v", err)
		return nil
	}

	// Unmarshell it to the struct
	err = json.Unmarshal([]byte(h.NotificationLogBytes), &h.NotificationLog)
	if err != nil {
		log.Println("error unmarshalling JSON: %v", err)
		return nil
	}
	// Marking day as true
	dum := make(map[int]bool)
	dum[time.Now().Day()] = true
	h.NotificationLog = dum

	h.NotificationLogBytes, err = json.Marshal(h.NotificationLog)
	if err != nil {
		return nil
	}
	return config.Rdb.HSet(context.Background(), key, "notification_log", h.NotificationLogBytes).Err()

}
func SetOffDay(key string, pipe redis.Pipeliner) redis.Pipeliner {
	h, err := getDaysRecord(key)
	if err != nil {
		log.Println("error getting days record: %v", err)
		return nil
	}

	// Unmarshell it to the struct
	err = json.Unmarshal([]byte(h.DaysLogByte), &h.DaysLog)
	if err != nil {
		log.Println("error unmarshalling JSON: %v", err)
		return nil
	}

	// Marking day as true
	h.DaysLog[time.Now().Day()] = false
	h.DaysLogByte, err = json.Marshal(h.DaysLog)
	if err != nil {
		return nil
	}
	pipe.HSet(context.Background(), key, "days_log", h.DaysLogByte)
	// New Joiner
	return pipe
}

// Since I would need to iterate over members multitimes, so why not making a multi use itrator!!
func Act(useCase string) {
	log.Println("Itratings")
	teleIDS, err := getMembersIDs()
	if err != nil {
		log.Println("err in InitOffDay while getting all member id: ", err)
		return
	}
	MembersCmdsMap := make(map[int]*redis.MapStringStringCmd)
	pipe := config.Rdb.Pipeline()
	for _, TId := range teleIDS {
		var key = RK(TId)
		if useCase == "SetDayOff" {
			pipe = SetOffDay(key, pipe)
		}

		MembersCmdsMap[TId] = pipe.HGetAll(context.Background(), key)

	}
	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Println("err pipe.Exec in InitOffDay while setting int day off for all members: ", err)
	}
	MemberActiveDaysMap := make(map[int]Habit)
	for teleID, cmd := range MembersCmdsMap {
		var h Habit
		err = cmd.Scan(&h)
		if err != nil {
			return
		}
		MemberActiveDaysMap[teleID] = h
	}
	if useCase == "SendStatus" {
		habitCalc(MemberActiveDaysMap)
	}
}

func habitCalc(memberActiveDaysMap map[int]Habit) {
	var highestStreakUser, highestTopHitUser Habit
	var highestStreak, highestTopHit int
	streakLeaderboard := []string{}

	for _, habit := range memberActiveDaysMap {
		// Track highest streak
		if habit.Streaked > highestStreak {
			highestStreak = habit.Streaked
			highestStreakUser = habit
		}
		// Track highest top hit (i.e., total number of days)
		if habit.TotalDays > highestTopHit {
			highestTopHit = habit.TotalDays
			highestTopHitUser = habit
		}
		// Add to streak leaderboard message (creative part)
		streakLeaderboard = append(streakLeaderboard, fmt.Sprintf(
			"🔥 %s is on a streak of %d days for habit **%s**!",
			habit.Name, habit.Streaked, habit.HabitName))
	}

	// Creative overall message
	topHitMsg := fmt.Sprintf(
		"🏅 Highest Top Hit: %s has completed **%d** days of habit **%s**! 🚀",
		highestTopHitUser.Name, highestTopHitUser.TotalDays, highestTopHitUser.HabitName)

	streakMsg := fmt.Sprintf(
		"🥇 Highest Streak: %s is on a **%d-day streak** for habit **%s**! Keep going! 🔥",
		highestStreakUser.Name, highestStreakUser.Streaked, highestStreakUser.HabitName)

	// Creative summary message
	summaryMsg := "📊 Daily Habit Overview:\n" +
		fmt.Sprintf("We’ve got some habit warriors making great progress today! 🌟\n") +
		strings.Join(streakLeaderboard, "\n") + "\n\n" +
		topHitMsg + "\n" + streakMsg

	// Send the message
	Remind(summaryMsg)
}
