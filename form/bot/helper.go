package bot

import (
	"familyFormUi/config"
	"fmt"
	"os"
	"strconv"

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

func Remind(text string) bool {
	botID, _ := strconv.Atoi(os.Getenv("TestingBotID"))
	config.B.Send(tele.ChatID(botID), text)
	return true
}
