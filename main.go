package main

import (
	"flag"
	"log"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	token := mustToken()

	tgClient := telegram.New(tgBotHost, token)

	_ = tgClient

	// TODO: make tg client

	// TODO: make fetcher

	// TODO: make processor

	// TODO: consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
