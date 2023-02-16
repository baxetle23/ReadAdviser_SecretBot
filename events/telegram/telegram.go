package telegram

import (
	"telegram_bot/clients"
	"telegram_bot/events"
	"telegram_bot/lib"
	"telegram_bot/storage"
)

type Processor struct {
	tg      *clients.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *clients.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, lib.WrapOnError("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func event(upd clients.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd clients.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd clients.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return lib.WrapOnError("can't process message", storage.ErrUnknownType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return lib.WrapOnError("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return lib.WrapOnError("can't process message", err)
	}

	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, lib.WrapOnError("can't get meta", storage.ErrUnknownType)
	}
	return res, nil
}
