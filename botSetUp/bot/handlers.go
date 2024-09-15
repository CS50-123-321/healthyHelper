package bot

import (
	"StreakHabitBulder/DB"
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
)

func (m *Members) Add(key string) (err error) {
	err = DB.Rdb.HSet(context.Background(), key, m).Err()
	if err != nil {
		return err
	}
	return
}

func getDaysRecord(key string) (h Habit, err error) {
	DB.Rdb.HGetAll(context.Background(), key).Scan(&h)
	return h, err
}

func (h *Habit) getStreakByUser() {

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
	err = DB.Rdb.ZRange(context.Background(), "MembersIDS", 0, -1).ScanSlice(&ids)
	if err != nil {
		return nil, err
	}
	return
}

// This get called by the cron job to run daily and sets the day as false, it will be true if the member did sport.
func InitOffDay() {
	teleIDS, err := getMembersIDs()
	if err != nil {
		log.Println("err in InitOffDay while getting all member id: ", err)
		return
	}
	pipe := DB.Rdb.Pipeline()
	for _, TId := range teleIDS {
		var h Habit
		var key = fmt.Sprintf("habitMember:%d", TId)
		h.DaysLogByte = []byte(fmt.Sprintf("{\"%v\":false}", strconv.Itoa(time.Now().Day())))
		pipe.HSet(context.Background(), key, "days_log", h.DaysLogByte)
	}
	Remind("Good Morning people, let's not forget to do sports today!!") // TODO: make better encourgment msg
	_, err = pipe.Exec(context.Background())
	if err != nil {
		log.Println("err pipe.Exec in InitOffDay while setting int day off for all members: ", err)
	}
}

// TODO: Do a daily report.
// TODO: Ranking, who is beating it.
// TODO: Who is missing out, who hasn't started.
