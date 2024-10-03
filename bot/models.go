package bot

import (
	"encoding/json"
	"fmt"
	"time"
)

func (m *Members) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Members) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

type Members struct {
	ID         int       `redis:"id" validate:"required"`
	Name       string    `redis:"name" validate:"required,min=2"`
	Habit      string    `redis:"habit" validate:"required,min=2"`
	StartDate  time.Time `redis:"start_date" validate:"required"`
	Days       []int     `redis:"days" validate:"required,min=1,dive,min=0,max=1"` // Checking all days as 0 (incomplete) or 1 (complete)
	TotalDays  int       `redis:"total_days" validate:"required,gt=0"`             // TotalDays must be greater than 0
	CurrentDay int       `redis:"current_day" validate:"required,gt=-1"`           // CurrentDay must be >= 0
}
type Streak struct {
	ID         int
	MemberID   int
	Habit      string
	StartDate  time.Time
	Days       []int
	TotalDays  int
	CurrentDay int
}

type Habit struct {
	Name                 string          `json:"name" redis:"name" binding:"required"`
	HabitName            string          `json:"habit_name" redis:"habit_name" binding:"required"`
	CommitmentPeriodStr  string          `json:"commitment_period"`
	CommitmentPeriod     int             `redis:"commitment_period"`
	TeleID               int             `redis:"tele_id"`
	TeleIDStr            string          `json:"tele_id"`
	Streaked             int             `redis:"streaked"`
	TopHit               int             `redis:"top_hit"` // the highest streak reached.
	DaysLog              map[string]bool // calc
	NotificationLog      map[string]bool
	NotificationLogBytes []byte `redis:"notification_log"`
	DaysLogByte          []byte `redis:"days_log"`
	TotalDays            int    `redis:"total_days"` // calc
	CreatedAt            time.Time
	GroupId              int `redis:"group_id" json:"group_id"`
}

func (h *Habit) MarshalBinary() ([]byte, error) {
	data, err := json.Marshal(h)
	if err != nil {
		return nil, fmt.Errorf("error marshalling habit to binary: %v", err)
	}
	return data, nil
}

func (h *Habit) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, h)
	if err != nil {
		return fmt.Errorf("error unmarshalling habit from binary: %v", err)
	}
	return nil
}
