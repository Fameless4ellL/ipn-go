package telegram

import (
	application "go-blocker/internal/application/payment"
	"go-blocker/internal/infrastructure/telegram/handlers"
	"go-blocker/internal/pkg/config"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Init(s *application.Service) {
	pref := tele.Settings{
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	hand := handlers.NewTGHandler(s)

	b.Handle("/healthcheck2", handlers.HealthCheck)
	b.Handle("/check", hand.CheckTx)
	b.Handle("/find", hand.FindTx)
	b.Handle("/SetChatID", handlers.SetChatID)
	go b.Start()
}
