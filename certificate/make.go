package certificate

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"time"

	"github.com/fogleman/gg"
)

type Habit struct {
	Name                 string       `json:"name" redis:"name" binding:"required"`
	HabitName            string       `json:"habitName" redis:"habit_name" binding:"required"`
	CommitmentPeriodStr  string       `json:"commitmentPeriod"`
	CommitmentPeriod     int          `redis:"commitment_period"`
	TeleID               int          `redis:"tele_id"`
	Streaked             int          `redis:"streaked"`
	TopHit               int          `redis:"top_hit"` // the highest streak reached.
	DaysLog              map[int]bool // calc
	NotificationLog      map[int]bool
	NotificationLogBytes []byte    `redis:"notification_log"`
	DaysLogByte          []byte    `redis:"days_log"`
	TotalDays            int       `redis:"total_days"` // calc
	CompletionDate       time.Time `redis:"completion_date"`
	// StartDate
	// Days count خابط
}

func GenerateCertificate(h Habit) (*bytes.Buffer, error) {
	const W = 900
	const H = 700

	dc := gg.NewContext(W, H)

	dc.SetRGB(0.95, 0.95, 1)
	dc.Clear()

	drawBackgroundPattern(dc, W, H)

	if err := dc.LoadFontFace("./Bangers-Regular.ttf", 48); err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}

	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawStringAnchored("🎉 Certificate of Achievement 🎉", W/2, 100, 0.5, 0.5)

	if err := dc.LoadFontFace("./Bangers-Regular.ttf", 36); err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}

	dc.SetRGB(0, 0, 0) // Black text
	message := fmt.Sprintf("This certifies that %s", h.Name)
	dc.DrawStringAnchored(message, W/2, 200, 0.5, 0.5)

	if err := dc.LoadFontFace("./Bangers-Regular.ttf", 42); err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}
	dc.SetRGB(0.8, 0.2, 0.2) // Bright red color
	habitMessage := fmt.Sprintf("has crushed the habit: %s 💪", h.HabitName)
	dc.DrawStringAnchored(habitMessage, W/2, 280, 0.5, 0.5)

	if err := dc.LoadFontFace("./Bangers-Regular.ttf", 28); err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}
	periodMessage := fmt.Sprintf("🗓️ Commitment Period: %d days", h.CommitmentPeriod)
	dateMessage := fmt.Sprintf("🏅 Completion Date: %s", h.CompletionDate.Format("January 02, 2006"))
	dc.DrawStringAnchored(periodMessage, W/2, 360, 0.5, 0.5)
	dc.DrawStringAnchored(dateMessage, W/2, 420, 0.5, 0.5)

	dc.SetLineWidth(6)
	dc.SetColor(color.RGBA{255, 215, 0, 255}) // Gold color for the border
	dc.DrawRectangle(30, 30, float64(W-60), float64(H-60))
	dc.Stroke()

	drawStarsAndConfetti(dc, W, H)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode certificate image: %v", err)
	}

	return buf, nil
}

func drawBackgroundPattern(dc *gg.Context, width, height int) {
	dc.SetLineWidth(2)
	dc.SetRGB(0.8, 0.9, 1)
	for i := 0; i < width; i += 20 {
		dc.DrawLine(float64(i), 0, float64(i), float64(height))
		dc.Stroke()
	}
}

func drawStarsAndConfetti(dc *gg.Context, width, height int) {
	dc.SetRGB(1, 0.8, 0)
	for i := 0; i < 15; i++ {
		dc.DrawLine(float64(i), 0, float64(i), float64(height))
		dc.Fill()
	}

	dc.SetRGB(0.2, 0.7, 0.4)
	for i := 0; i < 20; i++ {
		dc.DrawCircle(float64(100+i*40), float64(650+(i%2)*20), 5)
		dc.Fill()
	}
}
