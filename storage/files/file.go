package files

import (
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"os"
	"path/filepath"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	const op = "storage.files.Save"
	defer func() { err = e.WrapIfErr(op, err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath :=
	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

//Save(p *Page) error
//PickRandom(userName string) (*Page, error)
//Remove(p *Page) error
//IsExists(p *Page) (bool, error)
