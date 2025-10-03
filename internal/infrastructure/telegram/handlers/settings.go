package handlers

import (
	"fmt"
	"go-blocker/internal/pkg/config"

	tele "gopkg.in/telebot.v4"
)

func SetChatID(c tele.Context) error {
	chatID := c.Chat().ID
	// save at config
	config.ChatId = fmt.Sprint(chatID)
	return c.Send("Chat ID set to " + fmt.Sprint(chatID))
}
