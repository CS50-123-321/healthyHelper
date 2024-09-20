package bot

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// Daily reminder
// Daily ranking
// New joiners (say hi to x and y)
// Custom /comman for fetching the Top Ranking ETC
// Hero

func LevelMessage(h Habit, percentageCompleted int) error {
	var msg string
	var gifURL string

	switch {
	case percentageCompleted == 100:
		msg = fmt.Sprintf(
			"🏆 Congratulations, *%s*! You've completed *100%%* of your habit **%s**! You are officially a *Habit Champion*! 🎉 Keep up the great work and continue your streak! 💯🔥",
			h.Name, h.HabitName)
		gifURL = "https://vsgif.com/gif/3553131"
		// _, err := certificate.GenerateCertificate(h)
		// if err != nil {
		// 	return err
		// }
		// Config.B.
	case percentageCompleted == 80:
		msg = fmt.Sprintf(
			"💪 Amazing, *%s*! You've completed *60%%* of your habit **%s**! You're now a *Habit Hero*! Keep pushing forward, you're on fire! 🚀",
			h.Name, h.HabitName)
		gifURL = "https://media.giphy.com/media/xT9DPpf0zTqRASyzTi/giphy.gif"
	case percentageCompleted == 40:
		msg = fmt.Sprintf(
			"🌟 Great progress, *%s*! You've hit *30%%* of your habit **%s**! You're now a *Motivation Seeker*! Keep that momentum going! 💥",
			h.Name, h.HabitName)
		gifURL = "https://media.giphy.com/media/xT9IgG50Fb7Mi0prBC/giphy.gif"
	case percentageCompleted == 20:
		msg = fmt.Sprintf(
			"✨ Nice start, *%s*! You've completed *20%%* of your habit **%s**! You're officially a *Rising Star*! Keep up the effort, you've got this! ⭐",
			h.Name, h.HabitName)
		gifURL = "https://media.giphy.com/media/l46CjFkIMsxw6fQ5K/giphy.gif"
	case percentageCompleted == 0:
		msg = fmt.Sprintf(
			"🎉 Welcome, *%s*! 🎉\n\n"+
				"We're excited to have you on board for your new habit: **%s**! 💪\n"+
				"You've committed to building this habit for the next **%d days**. 🗓️\n\n"+
				"Stay strong and consistent, and we know you'll crush it! 🚀\n"+
				"Track your progress, stay motivated, and feel free to share your journey with the group! We're all cheering for you! 🙌✨\n\n",
			h.Name,
			h.HabitName,
			h.CommitmentPeriod,
		)
	}

	if msg != "" {
		// Setting the sent notigication true to avoid oversending msgs.
		err := SetNotificationLog(RK(h.TeleID))
		if err != nil {
			return err
		}
		// Send the message with bold formatting (MarkdownV2)
		err = Remind(msg)
		if err != nil {
			return err
		}
		botID, _ := strconv.Atoi(os.Getenv("TestingBotID"))
		// Send the GIF as an animation
		err = sendGIF(botID, os.Getenv("TELE_TOKEN"), gifURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendGIF(chatID int, botToken string, gifURL string) error {
	// Construct the Telegram API URL for sending the GIF
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendVideo?chat_id=%d&video=%s", botToken, chatID, gifURL)

	// Send the HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("error sending GIF: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send GIF, status code: %d", resp.StatusCode)
	}

	return nil
}
