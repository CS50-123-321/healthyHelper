package bot

import (
	"StreakHabitBulder/config"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	tele "gopkg.in/telebot.v3"
)

// Function to validate Members struct
func validateMembers(member *Members) error {
	validate := validator.New()
	err := validate.Struct(member)
	if err != nil {
		// Return validation errors
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Field '%s' failed validation, tag: '%s'\n", err.Field(), err.Tag())
		}
		return err
	}
	return nil
}

func Remind(text string) (err error) {
	log.Println("running remind")
	// botID, _ := strconv.Atoi(os.Getenv("TestingBotID"))
	// log.Println("----------------botID", botID)
	log.Println("rr", text)
	_, err = config.B.Send(tele.ChatID(-1002327721490), text, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2})
	if err != nil {
		log.Println("Remind: errsending the msg: ", err)
	}
	return err
}
func Personaliz(text, username string, userId int) (err error) {
	log.Println("Running personlize")
	text = FormatMention(username, userId)
	_, err = config.B.Send(tele.ChatID(-1002327721490), text, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, HasSpoiler: false})
	if err != nil {
		log.Println("Remind: errsending the msg: ", err)
	}
	return err
}
func RK(id int) string { return fmt.Sprintf("habitMember:%d", id) }
func FormatMention(Name string, teleID int) (msg string) {
	return fmt.Sprintf("[%s](tg://user?id=%d)", Name, teleID)
}
