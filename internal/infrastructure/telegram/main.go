package telegram

import (
	application "go-blocker/internal/application/payment"
	"go-blocker/internal/infrastructure/telegram/handlers"
	"go-blocker/internal/pkg/config"
	"log"
	"net/http"
	"net/url"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Init(s *application.Service) {
	prx, err := url.Parse("socks5://telegramik:telegramik@72.56.93.128:1080")
	if err != nil {
		panic(err)
	}

	pref := tele.Settings{
		Client: &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(prx)}},
		Token:  config.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Print(err)
		return
	}

	hand := handlers.NewTGHandler(s)

	b.Handle("/healthcheck2", handlers.HealthCheck)
	b.Handle("/check", hand.CheckTx)
	b.Handle("/find", hand.FindTx)
	b.Handle("/SetChatID", handlers.SetChatID)
	go b.Start()
}
