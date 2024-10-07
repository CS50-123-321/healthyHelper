package bot

import (
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
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

func GetDaysRecord(key string) (h Habit, err error) {
	config.Rdb.HGetAll(context.Background(), key).Scan(&h)
	return h, err
}

func (h *Habit) SetUserStreak() {

	h.TotalDays = 0
	h.Streaked = 0
	daysLogSlice := make([]string, 0, len(h.DaysLog))

	for dayStr := range h.DaysLog {
		daysLogSlice = append(daysLogSlice, dayStr)
	}
	sort.Strings(daysLogSlice)

	sortedMap := make(map[string]bool)
	for _, dayStr := range daysLogSlice {
		day, err := time.Parse("2006-01-02", dayStr)
		if err != nil {
			log.Println("err", err)
			return
		}
		if day.After(time.Now()) {
			return
		}
		done := h.DaysLog[dayStr]
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
		sortedMap[dayStr] = done
	}
	h.DaysLog = sortedMap
}

type MemberOrigin struct {
	TeleID, GroupId int
}

func getMembersIDs() (mo []MemberOrigin, err error) {
	membersWithScores, err := config.Rdb.ZRangeWithScores(context.Background(), "MembersIDS", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	for _, z := range membersWithScores {
		teleID, err := strconv.Atoi(z.Member.(string))
		if err != nil {
			return nil, fmt.Errorf("invalid TeleID: %v", err)
		}

		mo = append(mo, MemberOrigin{
			TeleID:  teleID,
			GroupId: int(z.Score),
		})
	}

	return mo, nil
}

func SetMemberLevel(memberHabit map[int]Habit) {
	for _, h := range memberHabit {
		percentageCompleted := 0
		if h.TotalDays > 1 {
			percentageCompleted = (h.TotalDays * 100 / h.CommitmentPeriod)
		}
		ok := h.NotificationLog[time.Now().Format("2006-01-02")]
		if !ok {
			LevelMessage(h, percentageCompleted)
		} else {
			log.Println("this user has been informed today already")
		}
	}
	return
}

// This get called by the cron job to run daily and sets the day as false, it will be true if the member did sport.

func SetNotificationLog(key string) error {
	h, err := GetDaysRecord(key)
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
	dum := make(map[string]bool)
	dum[time.Now().Format("2006-01-02")] = true
	h.NotificationLog = dum

	h.NotificationLogBytes, err = json.Marshal(h.NotificationLog)
	if err != nil {
		return nil
	}
	return config.Rdb.HSet(context.Background(), key, "notification_log", h.NotificationLogBytes).Err()

}

const (
	GenerateAiRandomMemberUseCASE = "GenerateAiRandomMember" // This sends ai genereted boost for random member in the goup.#AIGen
	MentionAllUseCASE             = "MentionAll"             // 	This sends a morning alike msg while mentioning everyone. #AIGen
	BestStreakUseCASE             = "bestStreak"             // This sends boost msg for the winner member. #AIGen
	DailyWatchUseCASE             = "dailyWatch"             //This sends positive msg to those who commit and sad song for those who didn't #AIGen
)

// Since I would need to iterate over members multitimes, so why not making a multi use itrator!!
func Act(useCase string) (habits []Habit) {
	err := Update()
	if err != nil {
		log.Println("err in update new version, redis shaping: ", err)
		return
	}
	log.Println("Itratings...")
	memberOrigin, err := getMembersIDs()
	if err != nil {
		log.Println("err in InitOffDay while getting all member id: ", err)
		return
	}
	if len(memberOrigin) == 0 {
		log.Println("getMembersIDs return empty slice")
		return []Habit{}
	}
	MembersCmdsMap := make(map[int]*redis.MapStringStringCmd)
	pipe := config.Rdb.Pipeline()
	for _, origin := range memberOrigin {
		var key = RK(origin.GroupId, origin.TeleID)
		MembersCmdsMap[origin.TeleID] = pipe.HGetAll(context.Background(), key)
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
		json.Unmarshal(h.DaysLogByte, &h.DaysLog)
		habits = append(habits, h)
		MemberActiveDaysMap[teleID] = h
	}

	switch useCase {
	case DailyWatchUseCASE:
		log.Println("dailyWatch..")
		DailyWatch(MemberActiveDaysMap)
	case BestStreakUseCASE:
		log.Println("bestStreak..")
		BestStreak(habits)
	case MentionAllUseCASE:
		log.Println("MentionAll..")
		MentionAll(habits)
	case GenerateAiRandomMemberUseCASE:
		log.Println("SendAiPersonalizedMsg..")
		SendAiPersonalizedMsg(habits)
	}
	return habits
}

// GetHabitLevel calculates the habit level based on the commitment period and days completed
func GetHabitLevel(completionPercentage int) string {
	// Determine the level based on the completion percentage
	switch {
	case completionPercentage >= 100:
		return "Habit Hero 🏆"
	case completionPercentage >= 40:
		return "Motivation Seeker 🚀"
	case completionPercentage >= 20:
		return "Rising Star 🌟"
	default:
		return "New Challenger🌱" // 0% level
	}
}

type Tag struct {
	TagBody string
	Streak  int
}

func Update() error {
	ids := []int{}
	err := config.Rdb.ZRange(context.Background(), "MembersIDS", 0, -1).ScanSlice(&ids)
	if err != nil {
		return err
	}
	for _, id := range ids {
		var h Habit
		err = config.Rdb.ZAdd(context.Background(), "MembersIDS", redis.Z{
			Score:  -1002327721490,
			Member: id,
		}).Err()
		if err != nil {
			return err
		}
		err = config.Rdb.HGetAll(context.Background(), fmt.Sprintf("habitMember:%v", id)).Scan(&h)
		if err != nil {
			return err
		}
		h.GroupId = -1002327721490
		var key = RK(h.GroupId, id)
		err = config.Rdb.HSet(context.Background(), key, h).Err()
		if err != nil {
			return err
		}
		err = config.Rdb.SAdd(context.Background(), "groupIds", h.GroupId).Err()
		if err != nil {
			return err
		}
		err = config.Rdb.SAdd(context.Background(), fmt.Sprintf("habitByGroup:%v", h.GroupId), id).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
