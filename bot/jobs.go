package bot

import (
	"StreakHabitBulder/config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"math/rand"
)

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
		Remind(msg)
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
		Remind(fmt.Sprintf("%s\n%s", msg, EscapeMarkdown(MentionAllBody)))
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
	Remind(fmt.Sprintf("%s \n %s", ExecAbleBody, AiResponse))
}

func DailyWatch(memberActiveDaysMap map[int]Habit) {
	var p HabitMessage
	for _, h := range memberActiveDaysMap {
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
		if msg != "" {
			Remind(EscapeMarkdown(msg), tag)
		}
	}
}
