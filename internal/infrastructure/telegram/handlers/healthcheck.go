package handlers

import (
	tele "gopkg.in/telebot.v4"
)

func HealthCheck(c tele.Context) error {
	return c.Send("Bot is running!")
}
