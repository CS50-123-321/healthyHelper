package bot

import (
	"fmt"
)

type HabitMessage struct {
	Base struct {
		AllowedSymobls string
	}
	DailyWatch struct {
		Committed   string
		NotCommited string
	}
	InstantReply struct {
		AfterDayCounter string
	}
}

// this helps in the dynamic fetching without the need to prepre all promots, on need approach instead.
const (
	Dailywtch = "DailyWatch"
)

func (hs *HabitMessage) HabitMsgs(h Habit, category string) {
	hs.Base.AllowedSymobls = "USE ONLY THE FOLLOWING SYMBOLS IN THE TEXT:  '!',',','?','#'"
	hs.InstantReply.AfterDayCounter = "In Arabic, Generate five to 10 words sentence simple Arabic poem that is easy to understand for young adults aged 15-22. The poem should convey themes of determination, positivity, and encouragement. Use creative and inspiring Arabic language that avoids clichés and motivates action and perseverance., don't translate"
	switch category {
	case Dailywtch:
		hs.DailyWatch.Committed = fmt.Sprintf("%s has committed to do %s for today, with today he has been doing his habit for %v. Generte an encouring one line msg for them. %s", h.Name, h.HabitName, h.TotalDays, hs.Base.AllowedSymobls)
		hs.DailyWatch.NotCommited = fmt.Sprintf("Write a lighthearted, 40-word song in Arabic about a young person named %s who missed doing their habit called %s today. Make it fun and relatable for 12-19-year-olds, using emojis and humor to remind them of the importance of staying committed to their habit, but in a playful way that encourages them to try again tomorrow. %s", h.Name, h.HabitName, hs.Base.AllowedSymobls)
	}
}
