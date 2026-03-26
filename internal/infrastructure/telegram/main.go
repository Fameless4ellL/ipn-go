package telegram

import (
	"context"
	application "go-blocker/internal/application/payment"
	"go-blocker/internal/infrastructure/telegram/handlers"
	"go-blocker/internal/pkg/config"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
	tele "gopkg.in/telebot.v4"
)

func Init(s *application.Service) {
	proxyAddr := "72.56.93.128:1080"
	auth := &proxy.Auth{
		User:     "telegramik",
		Password: "telegramik",
	}

	// 2. Create a SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, auth, proxy.Direct)
	if err != nil {
		panic(err)
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}

	transport := &http.Transport{
		DialContext: dialContext,
	}

	proxyClient := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 30,
	}

	pref := tele.Settings{
		Client: proxyClient,
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
