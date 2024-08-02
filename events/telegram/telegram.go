package telegram

import (
	"errors"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/clients/telegram"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/events"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/lib/e"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) Processor {
	return Processor{
		tg:      client,
		offset:  100,
		storage: storage,
	}
}

func (c *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := c.tg.Updates(c.offset, limit)

	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return []events.Event{}, nil
	}

	res := make([]events.Event, 0, len(updates))

	for i, u := range updates {
		res[i] = event(u)
	}

	c.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	defer func() { err = e.Wrap("can't process message", err) }()

	meta, err := meta(event)

	if err != nil {
		return err
	}

	return p.doCmd(event.Text, meta.ChatId, meta.Username)
}

func event(u telegram.Update) events.Event {
	updType := fetchType(u)

	res := events.Event{
		Type: updType,
		Text: fetchText(u),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func meta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return meta, nil
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}

	return u.Message.Text
}
