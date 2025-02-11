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

func Remind(text string, chatID int, tag ...string) (err error) {
	log.Println("running remind")
	if len(tag) == 1 {
		text = fmt.Sprintf("%s \n %s", tag[0], text)
	}
	if chatID == 0 {
		chatID, _ = strconv.Atoi(os.Getenv("TestingBotID"))

	}
	_, err = config.B.Send(tele.ChatID(chatID), text, &tele.SendOptions{ParseMode: tele.ModeMarkdownV2, HasSpoiler: false})
	if err != nil {
		log.Println("Remind: errsending the msg: ", err, text)
	}
	return err
}
func RK(groupId, id int) string { return fmt.Sprintf("habitMember:%v:%v", groupId, id) }
func FormatMention(Name string, teleID int) (msg string) {
	return fmt.Sprintf("[%s](tg://user?id=%d)", Name, teleID)
}

func EscapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}
