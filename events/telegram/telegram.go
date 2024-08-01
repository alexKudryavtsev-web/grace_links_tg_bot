package telegram

import "github.com/alexKudryavtsev-web/grace_links_tg_bot/clients/telegram"

type Consumer struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(client *telegram.Client) Consumer {
	return Consumer{
		tg:     client,
		offset: 100,
	}
}
