package bot

import (
	"fmt"
	"log"
	"sort"
	"time"
)

func BestStreak(AllMemberHabits []Habit) {
	if len(AllMemberHabits) == 0 {
		log.Println("BestStreak, no members to get the best of them")
		return
	}
	sort.Slice(AllMemberHabits, func(i, j int) bool {
		return AllMemberHabits[i].TotalDays > AllMemberHabits[j].TotalDays
	})
	filteredMembers := []Habit{}
	for _, habit := range AllMemberHabits {
		if habit.TotalDays > 1 {
			filteredMembers = append(filteredMembers, habit)
		}
	}

	if len(filteredMembers) == 0 {
		log.Println("BestStreak, no one has done anything impressive so far, keep pushing!")
		return
	}

	numTopUsers := len(filteredMembers)
	medalEmojis := []string{"🥇", "🥈", "🥉"}

	switch numTopUsers {
	case 1:
		// For one user
		progressMsg := fmt.Sprintf("🏆 *Today's Winner:* 🏆\n\n🥇 %s: Total Days: %d, Streak: %d days, Habit: %s\n\n",
			FormatMention(filteredMembers[0].Name, filteredMembers[0].TeleID),
			filteredMembers[0].TotalDays,
			filteredMembers[0].Streaked,
			filteredMembers[0].HabitName)
		Remind(progressMsg, 0)
	case 2:
		// For two users
		progressMsg := "🏆 *Today's Top 2 Winners:* 🏆\n\n"
		for i := 0; i < 2; i++ {
			progressMsg += fmt.Sprintf("%s %s: Total Days: %d, Streak: %d days, Habit: %s\n\n",
				medalEmojis[i],
				FormatMention(filteredMembers[i].Name, filteredMembers[i].TeleID),
				filteredMembers[i].TotalDays,
				filteredMembers[i].Streaked,
				filteredMembers[i].HabitName)
		}
		Remind(progressMsg, 0)
	case 3:
		progressMsg := "🏆 *Today's Top 3 Winners:* 🏆\n\n"
		for i := 0; i < 3; i++ {
			progressMsg += fmt.Sprintf("%s  %s: Total Days: %d, Streak: %d days, Habit: %s\n\n",
				medalEmojis[i],
				FormatMention(filteredMembers[i].Name, filteredMembers[i].TeleID),
				filteredMembers[i].TotalDays,
				filteredMembers[i].Streaked,
				filteredMembers[i].HabitName)
		}
		Remind(progressMsg, 0)
	case 4:
		progressMsg := "🏆 *Today's Top 4 Winners:* 🏆\n\n"
		for i := 0; i < 4; i++ {
			emoji := medalEmojis[i%3]
			if i >= 3 {
				emoji = "⭐️"
			}
			progressMsg += fmt.Sprintf("%s  %s: Total Days: %d, Streak: %d days, Habit: %s\n\n",
				emoji,
				FormatMention(filteredMembers[i].Name, filteredMembers[i].TeleID),
				filteredMembers[i].TotalDays,
				filteredMembers[i].Streaked,
				filteredMembers[i].HabitName)
		}
		Remind(progressMsg, 0)
	case 5:
		progressMsg := "🏆 *Today's Top 5 Winners:* 🏆\n\n"
		for i := 0; i < 5; i++ {
			emoji := medalEmojis[i%3]
			if i >= 3 {
				emoji = "⭐️"
			}
			progressMsg += fmt.Sprintf("%s  %s: Total Days: %d, Streak: %d days, Habit: %s\n\n",
				emoji,
				FormatMention(filteredMembers[i].Name, filteredMembers[i].TeleID),
				filteredMembers[i].TotalDays,
				filteredMembers[i].Streaked,
				filteredMembers[i].HabitName)
		}
		Remind(progressMsg, 0)
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
		Remind(fmt.Sprintf("%s\n%s", msg, MentionAllBody), 0)
	}
}

func SendAiPersonalizedMsg(habits []Habit) {
	// var tries int = len(habits)
	// if tries == 0 { // to avoid infinite recursivnessnessnessnessness
	// 	log.Println("all members have recieved an ai generated boost")
	// 	return
	// }
	// rndIndx := rand.Intn(len(habits))
	// h := habits[rndIndx]
	// found := config.Rdb.SIsMember(context.Background(), "sentList:membersIDS", h.TeleID).Val()
	// if found {
	// 	habits = append(habits[:rndIndx], habits[rndIndx+1:]...) // cutting the already sent to member and procceed.
	// 	SendAiPersonalizedMsg(habits)                            // recursively calling back until reaching new memeber
	// 	return
	// }
	// config.Rdb.SAdd(context.Background(), "sentList:membersIDS", h.TeleID)
	// now := time.Now()
	// midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	// durationUntilMidnight := time.Until(midnight) //set it tobe deletedd after the end of todayyyyy
	// err := config.Rdb.Conn().Expire(context.Background(), "sentList:membersIDS", durationUntilMidnight).Err()
	// if err != nil {
	// 	fmt.Println("Error setting TTL:", err)
	// }
	for _, h := range habits {
		AiResponse, err := GetAiResponse(h)
		if err != nil {
			log.Println("err SendAiPersonalizedMsg", err)
			return
		}
		AiResponse = EscapeMarkdown(AiResponse)
		tag := FormatMention(h.Name, h.TeleID)
		Remind(AiResponse, h.TeleID, tag)
		time.Sleep(10 * time.Second)
	}

}

func DailyWatch(memberActiveDaysMap map[int]Habit) {
	var p HabitMessage
	for _, h := range memberActiveDaysMap {
		var msg string
		var err error
		done, ok := h.DaysLog[time.Now().Day()]
		tag := FormatMention(h.Name, h.TeleID)
		p.HabitMsgs(h, Dailywtch)
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
			Remind(EscapeMarkdown(msg), 0, tag)
		}
	}
}
