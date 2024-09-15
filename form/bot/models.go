package bot

import (
	"encoding/json"
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
