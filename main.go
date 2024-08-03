package main

import (
	"flag"
	"log"

	tgClient "github.com/alexKudryavtsev-web/grace_links_tg_bot/clients/telegram"
	event_consumer "github.com/alexKudryavtsev-web/grace_links_tg_bot/consumer/event-consumer"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/storage/files"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/events/telegram"
)

const (
	tgBotHost = "api.telegram.org"
	storage   = "file-storage"
	batchSize = 10
)

func main() {
	token := mustToken()
	tgClient := tgClient.New(tgBotHost, token)
	processor := telegram.New(&tgClient, files.New(storage))
 
	log.Println("server started")

	consumer := event_consumer.New(&processor, &processor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
