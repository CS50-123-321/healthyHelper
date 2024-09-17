package api

type Habit struct {
	Name                 string       `json:"name" redis:"name" binding:"required"`
	HabitName            string       `json:"habitName" redis:"habit_name" binding:"required"`
	CommitmentPeriodStr  string       `json:"commitmentPeriodStr" redis:"commitmentPeriodStr"`
	CommitmentPeriod     int          `redis:"commitment_Period"`
	TeleID               int64        `redis:"tele_id"`
	Streaked             int64        `redis:"streaked"`
	DaysLog              map[int]bool // calc
	DaysLogByte          []byte       `redis:"days_log"`
	NotificationLog      map[int]bool
	NotificationLogBytes []byte `redis:"notification_log"`
	TotalDays            int64  `redis:"total_days"` // calc
	TopHit               int    `redis:"top_hit"`   // the highest dtreak reached.
}
