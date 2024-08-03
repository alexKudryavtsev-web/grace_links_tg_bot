package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/alexKudryavtsev-web/grace_links_tg_bot/clients/telegram"
	event_consumer "github.com/alexKudryavtsev-web/grace_links_tg_bot/consumer/event-consumer"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/storage/sqlite"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/events/telegram"
)

const (
	tgBotHost = "api.telegram.org"
	storage   = "data/sqlite/storage.sqlite"
	batchSize = 100
)

func main() {
	token := mustToken()
	tgClient := tgClient.New(tgBotHost, token)

	db, err := sqlite.New(storage)

	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	if err := db.Init(context.TODO()); err != nil {
		log.Fatalf("can't init storage")
	}

	processor := telegram.New(&tgClient, db)

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
