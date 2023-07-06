package storage

import (
	"TelegramBot/lib/e"
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (res string, err error) {
	const op = "storage.Hash"
	defer func() { err = e.WrapIfErr(op, err) }()

	h := sha1.New()

	if _, err = io.WriteString(h, p.URL); err != nil {
		return res, err
	}

	if _, err = io.WriteString(h, p.UserName); err != nil {
		return res, err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
