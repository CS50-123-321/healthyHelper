package bot

import (
	"fmt"
	"net/http"
	"os"
)

func LevelMessage(h Habit, percentageCompleted int) (err error) {
	var msg string
	var gifURL string
	switch {
	case percentageCompleted == 100:
		msg = fmt.Sprintf(
			"🏆 Congratulations, *%s*! You've completed *100%%* of your habit **%s**! You are officially a *Habit Champion*! 🎉 Keep up the great work and continue your streak! 💯🔥",
			h.Name, h.HabitName)

		gifURL = "https://media1.giphy.com/media/UZOny91tZDtDyXPrJC/giphy.gif?cid=9b38fe91mnvzk7h19m15st9wp5fuog8o7n4l6u85rzcx0lnm&ep=v1_gifs_search&rid=giphy.gif&ct=g"
		//gold star for you git https://media3.giphy.com/media/LLYtaSxAuj47fmmwpU/giphy.gif?cid=9b38fe91ye663gwag3ojjunzdw9kk6akwygh47xmz7ttq24l&ep=v1_gifs_search&rid=giphy.gif&ct=g
		// _, err := certificate.GenerateCertificate(h)
		// if err != nil {
		// 	return err
		// }
		// Config.B.
	case percentageCompleted == 80:
		msg = fmt.Sprintf(
			"💪 Amazing, *%s*! You've completed *60%%* of your habit **%s**! You're now a *Habit Hero*! Keep pushing forward, you're on fire! 🚀",
			h.Name, h.HabitName)
		gifURL = "https://media2.giphy.com/media/jQhdmWUmKgn9bZJxUj/giphy.gif?cid=9b38fe916n4x3ak9bjefcci615suljtzn71m0woqv7unk3zl&ep=v1_gifs_search&rid=giphy.gif&ct=g"

	case percentageCompleted == 40:
		msg = fmt.Sprintf(
			"🌟 Great progress, *%s*! You've hit *30%%* of your habit **%s**! You're now a *Motivation Seeker*! Keep that momentum going! 💥",
			h.Name, h.HabitName)
		gifURL = "https://media0.giphy.com/media/nhRDKEdPXWaA0NZWxT/giphy.gif?cid=9b38fe91mtvmxvklctoczirh4ig46r0sqenmsiokykd75142&ep=v1_gifs_search&rid=giphy.gif&ct=g"
	case percentageCompleted == 20:
		msg = fmt.Sprintf(
			"✨ Nice start, *%s*! You've completed *20%%* of your habit **%s**! You're officially a *Rising Star*! Keep up the effort, you've got this! ⭐",
			h.Name, h.HabitName)
		gifURL = "https://media0.giphy.com/media/xUNd9Yd1YNOiFd2D16/giphy.gif?cid=9b38fe912hu6qe4ovfmnivjj38rm2175wa2e9jktd4f7i6zn&ep=v1_gifs_search&rid=giphy.gif&ct=g"
	case percentageCompleted == 0:
		msg = fmt.Sprintf(
			"🎉 Welcome, *%s*! 🎉\n\n"+
				"We're excited to have you on board for your new habit: **%s**! 💪\n"+
				"You've committed to building this habit for the next **%d days**. 🗓️\n\n"+
				"Stay strong and consistent, and we know you'll crush it! 🚀\n"+
				"Track your progress, stay motivated, and feel free to share your journey with the group! We're all cheering for you! 🙌✨\n\n"+
				// Adding a space between the English and Arabic sections
				"\n\n"+
				"🎉 مرحبًا، *%s*! 🎉\n\n"+
				"يسعدنا انضمامك إلى عادتك الجديدة: **%s**! 💪\n"+
				"لقد التزمت ببناء هذه العادة خلال **%d من الأيام القادمة**. 🗓\n\n"+
				"ابق قويًا وثابتًا، ونحن نعلم أنك ستفوز! 🚀\n"+
				"تابع تقدمك، وحافظ على تحفيزك، ولا تتردد في مشاركة رحلتك مع المجموعة! نحن جميعًا نشجعك! 🙌✨\n\n",
			h.Name,
			h.HabitName,
			h.CommitmentPeriod,
			h.Name,
			h.HabitName,
			h.CommitmentPeriod,
		)
		gifURL = "https://media1.giphy.com/media/443kbHqcFFwOvTuAWf/giphy.gif?cid=9b38fe916n4x3ak9bjefcci615suljtzn71m0woqv7unk3zl&ep=v1_gifs_search&rid=giphy.gif&ct=g"

	}
	if msg != "" {
		// Setting the sent notigication true to avoid oversending msgs.
		err := SetNotificationLog(RK(h.GroupId, h.TeleID))
		if err != nil {
			return err
		}
		// Send the message with bold formatting (MarkdownV2)
		err = Remind(EscapeMarkdown(msg), h.GroupId)
		if err != nil {
			return err
		}
		// Send the GIF as an animation
		err = sendGIF(h.GroupId, os.Getenv("TELE_TOKEN"), gifURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendGIF(groupID int, botToken string, gifURL string) error {
	// Construct the Telegram API URL for sending the GIF
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendVideo?chat_id=%d&video=%s", botToken, groupID, gifURL)

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
