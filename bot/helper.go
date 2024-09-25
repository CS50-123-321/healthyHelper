package bot

import (
	"StreakHabitBulder/config"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	botID, _ := strconv.Atoi(os.Getenv("TestingBotID"))
	log.Println("botid", botID)
	_, err = config.B.Send(tele.ChatID(-4580179828), text, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, HasSpoiler: false})
	if err != nil {
		log.Println("Remind: errsending the msg: ", err)
	}
	return err
}
func RK(id int) string { return fmt.Sprintf("habitMember:%d", id) }
func FormatMention(Name string, teleID int) (msg string) {
	return fmt.Sprintf("[%s](tg://user?id=%d)", Name, teleID)
}

func EscapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"-", "\\-",
		"!", "\\!",
		"#", "\\",
		".", "\\.",
	)
	return replacer.Replace(text)
}
