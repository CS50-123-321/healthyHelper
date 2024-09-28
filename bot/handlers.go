package bot

import (
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

func GetDaysRecord(key string) (h Habit, err error) {
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
		if day > time.Now().Day() {
			return
		}
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
		if h.TotalDays > 1 {
			percentageCompleted = (h.TotalDays * 100 / h.CommitmentPeriod)
		}
		ok := h.NotificationLog[time.Now().Day()]
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
	h, err := GetDaysRecord(key)
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

const (
	GenerateAiRandomMemberUseCASE = "GenerateAiRandomMember" // This sends ai genereted boost for random member in the goup.#AIGen
	MentionAllUseCASE             = "MentionAll"             // 	This sends a morning alike msg while mentioning everyone. #AIGen
	BestStreakUseCASE             = "bestStreak"             // This sends boost msg for the winner member. #AIGen
	DailyWatchUseCASE             = "dailyWatch"             //This sends positive msg to those who commit and sad song for those who didn't #AIGen
)

// Since I would need to iterate over members multitimes, so why not making a multi use itrator!!
func Act(useCase string) (habits []Habit) {
	log.Println("Itratings...")
	teleIDS, err := getMembersIDs()
	if err != nil {
		log.Println("err in InitOffDay while getting all member id: ", err)
		return
	}
	if len(teleIDS) == 0 {
		log.Println("getMembersIDs return empty slice")
		return []Habit{}
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
		habits = append(habits, h)
		MemberActiveDaysMap[teleID] = h
	}

	switch useCase {
	case "SendStatus":
		log.Println("SendStatus..")
		habitCalc(MemberActiveDaysMap)
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

func habitCalc(memberActiveDaysMap map[int]Habit) {
	var highestStreakUser, highestTopHitUser Habit
	var highestStreak, highestTopHit int
	streakLeaderboard := []string{}

	for _, habit := range memberActiveDaysMap {
		if habit.Streaked > highestStreak {
			highestStreak = habit.Streaked
			highestStreakUser = habit
		}
		if habit.TotalDays > highestTopHit {
			highestTopHit = habit.TotalDays
			highestTopHitUser = habit
		}
		streakLeaderboard = append(streakLeaderboard, fmt.Sprintf(
			"🔥 %s is on a streak of %d days for habit **%s**",
			FormatMention(habit.Name, habit.TeleID), habit.Streaked, habit.HabitName))
	}

	topHitMsg := fmt.Sprintf(
		"🏅 Highest Top Hit: %s has completed **%d** days of habit **%s** 🚀",
		FormatMention(highestTopHitUser.Name, highestTopHitUser.TeleID), highestTopHitUser.TotalDays, highestTopHitUser.HabitName)

	streakMsg := fmt.Sprintf(
		"🥇 Highest Streak: %s is on a **%d day streak** for habit **%s** Keep going 🔥",
		FormatMention(highestTopHitUser.Name, highestTopHitUser.TeleID), highestStreakUser.Streaked, highestStreakUser.HabitName)

	summaryMsg := "📊 Daily Habit Overview:\n" +
		fmt.Sprintf("We’ve got some habit warriors making great progress today 🌟\n") +
		strings.Join(streakLeaderboard, "\n") + "\n\n" +
		topHitMsg + "\n" + streakMsg
	chatID := 0
	// if highestStreakUser.IsGroup {
	// 	chatID = highestStreakUser.TeleID
	// } this is for the group only
	Remind(summaryMsg, chatID)
}

func DailyWatch(memberActiveDaysMap map[int]Habit) {
	var p HabitMessage
	for _, h := range memberActiveDaysMap {
		log.Println("00000000", h.Name, h.IsGroup)

		err := json.Unmarshal(h.DaysLogByte, &h.DaysLog)
		if err != nil {
			log.Println(err)
			return
		}
		var msg string
		done, ok := h.DaysLog[time.Now().Day()]
		tag := FormatMention(h.Name, h.TeleID)
		p.HabitMsgs(h, Dailywtch) // this filles the structs the promits based on the function need.
		if ok && !done {
			msg, err = GenerateText(p.DailyWatch.NotCommited)
			if err != nil {
				log.Println(err)
				return
			}
		} else if ok && done {
			msg, err = GenerateText(p.DailyWatch.Committed)
			if err != nil {
				log.Println(err)
				return
			}
		}
		chatID := 0
		if !h.IsGroup {
			chatID = h.TeleID
		}
		if msg != "" {
			Remind(EscapeMarkdown(msg), chatID, tag)
		}
	}
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

func BestStreak(AllMemberHabits []Habit) {
	if len(AllMemberHabits) == 0 {
		log.Println("BestStreak, no members to get the best of them")
		return
	}
	sort.Slice(AllMemberHabits, func(i, j int) bool {
		return AllMemberHabits[i].TotalDays > AllMemberHabits[j].TotalDays
	})
	topDays := AllMemberHabits[0].TotalDays
	if topDays == 0 {
		log.Println("BestStreak, no one has done anything impressive sofar, fuck off")
		return
	}
	topUsers := []Tag{}
	for _, habit := range AllMemberHabits {
		if habit.TotalDays == topDays {
			topUsers = append(topUsers, Tag{
				TagBody: FormatMention(habit.Name, habit.TeleID),
				Streak:  habit.Streaked,
			})
		} else {
			break
		}
	}
	for _, tag := range topUsers {
		msg := fmt.Sprintf("Look at you go\\!\\! \n %s You're already at %v days\\. One step closer to being a habit hero\\!", tag.TagBody, tag.Streak)
		Remind(msg, 0) //TODO: make this sending to the user dm if solo learner.
	}
}

var maxRetriesLimit int = 3

func MentionAll(habits []Habit) {
	var promptLanguage []string = []string{"In English, Generate a morning message for group of habit builders, it has to be cool and motivating",
		"In Arabic, Generate a morning message for group of habit builders, it has to be cool and motivating"}
	var MentionAllBody string
	for _, h := range habits {
		MentionAllBody = MentionAllBody + fmt.Sprintf(" %s, ", FormatMention(h.Name, h.TeleID))
	}
	for i := range promptLanguage {
		p := promptLanguage[i]
		AiResponse, err := GenerateText(p)
		if maxRetriesLimit == 0 {
			log.Println("reaching max tried in mentionAll")
			return
		}
		if err != nil {
			maxRetriesLimit--
			GenerateText(p)
			return
		}
		msg := EscapeMarkdown(AiResponse)
		Remind(fmt.Sprintf("%s\n%s", msg, EscapeMarkdown(MentionAllBody)), 0) // this doesn't work for now. TODO:
	}
}

func SendAiPersonalizedMsg(habits []Habit) {
	var tries int = len(habits)
	if tries == 0 { // to avoid infinite recursivnessnessnessnessness
		log.Println("all members have recieved an ai generated boost")
		return
	}
	rndIndx := rand.Intn(len(habits))
	h := habits[rndIndx]
	found := config.Rdb.SIsMember(context.Background(), "sentList:membersIDS", h.TeleID).Val()
	if found {
		habits = append(habits[:rndIndx], habits[rndIndx+1:]...) // cutting the already sent to member and procceed.
		SendAiPersonalizedMsg(habits)                            // recursively calling back until reaching new memeber
		return
	}
	config.Rdb.SAdd(context.Background(), "sentList:membersIDS", h.TeleID)
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	durationUntilMidnight := time.Until(midnight) //set it tobe deletedd after the end of todayyyyy
	err := config.Rdb.Conn().Expire(context.Background(), "sentList:membersIDS", durationUntilMidnight).Err()
	if err != nil {
		fmt.Println("Error setting TTL:", err)
	}

	AiResponse, err := GetAiResponse(h)
	if err != nil {
		log.Println("err SendAiPersonalizedMsg", err)
		return
	}
	AiResponse = EscapeMarkdown(AiResponse)
	ExecAbleBody := FormatMention(h.Name, h.TeleID)
	chatID := 0
	if h.IsGroup {
		chatID = h.TeleID
	}
	Remind(fmt.Sprintf("%s \n %s", ExecAbleBody, AiResponse), chatID)
}
