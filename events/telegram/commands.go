package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"TelegramBot/clients/telegram"
	"TelegramBot/lib/e"
	"TelegramBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL, username string) (err error) {
	const op = "events.telegram.savePage"
	defer func() { err = e.WrapIfErr(op, err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExist {
		return sendMsg(msgAlreadyExists)
	}

	if err = p.storage.Save(page); err != nil {
		return err
	}

	if err = p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	const op = "events.telegram.savePage"
	defer func() { err = e.WrapIfErr(op, err) }()
	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavePages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavePages) {
		return sendMsg(msgNoSavedPages)
	}

	if err = sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) sendHello(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// http://ya.ru
	u, err := url.Parse(text)

	if err == nil && u.Host != "" {
		return true
	}
	return false
}
