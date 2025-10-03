package telegram

import (
	"go-blocker/internal/infrastructure/telegram/handlers"
	"go-blocker/internal/pkg/config"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Init() {
	pref := tele.Settings{
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// db := database.New()
	// repo := database.NewPaymentRepository(db)
	// service := payment.NewPaymentService(repo)

	b.Handle("/healthcheck", handlers.HealthCheck)
	b.Handle("/checktx", handlers.CheckTx)
	b.Handle("/findtx", handlers.FindTx)
	b.Start()
}
