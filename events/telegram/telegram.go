package telegram

import (
	"TelegramBot/clients/telegram"
	"TelegramBot/events"
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"errors"
)

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	const op = "events.telegram.Fetch"

	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap(op, err)
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

func (p *Processor) Process(event events.Event) error {
	const op = "events.telegram.Process"

	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap(op, ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) (err error) {
	const op = "events.telegram.processMessage"
	defer func() { err = e.WrapIfErr(op, err) }()

	meta, err := meta(event)
	if err != nil {
		return err
	}
	if err = p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	const op = "events.telegram.meta"

	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap(op, ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	// chatID username
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res

}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
