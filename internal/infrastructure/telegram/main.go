package telegram

import (
	application "go-blocker/internal/application/payment"
	"go-blocker/internal/infrastructure/telegram/handlers"
	"go-blocker/internal/pkg/config"
	"log"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v4"
)

type HeaderTransport struct {
	Header http.Header
	Base   http.RoundTripper
}

func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, vv := range t.Header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	return t.Base.RoundTrip(req)
}

func Init(s *application.Service) {
	header := http.Header{}
	header.Add("authorization", config.TELEGRAM_AUTH_BASE_URL)

	customClient := &http.Client{
		Transport: &HeaderTransport{
			Header: header,
			Base:   http.DefaultTransport,
		},
		Timeout: time.Minute,
	}

	pref := tele.Settings{
		URL:   config.TG_BASE_URL,
		Token: config.BotToken,
		Poller: &tele.Webhook{
			Endpoint:         &tele.WebhookEndpoint{PublicURL: config.TG_WEBHOOK_PATH},
			Listen:           "localhost:8888",
			IgnoreSetWebhook: true,
		},
		Client: customClient,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	hand := handlers.NewTGHandler(s)

	b.Handle("/healthcheck", handlers.HealthCheck)
	b.Handle("/check", hand.CheckTx)
	b.Handle("/find", hand.FindTx)
	b.Handle("/SetChatID", handlers.SetChatID)
	go b.Start()
}
