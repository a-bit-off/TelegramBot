package files

import (
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	const op = "storage.files.Save"
	defer func() { err = e.WrapIfErr(op, err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (p *storage.Page, err error) {
	const op = "storage.files.PickRandom"
	defer func() { err = e.WrapIfErr(op, err) }()

	fPath := filepath.Join(s.basePath, userName)

	// TODO: refactor: PickRandom
	// 1. check user folder
	// 2. create folder

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	filesLen := len(files)
	if filesLen == 0 {
		return nil, storage.ErrNoSavePages
	}

	//	0 - (len - 1)
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(filesLen)

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) (err error) {
	const op = "storage.files.Remove"
	defer func() { err = e.WrapIfErr(op, err) }()

	fName, err := fileName(p)
	if err != nil {
		return err
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err = os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (success bool, err error) {
	const op = "storage.files.IsExists"
	defer func() { err = e.WrapIfErr(op, err) }()

	fName, err := fileName(p)
	if err != nil {
		return false, err
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	const op = "storage.files.decodePage"
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err = gob.NewDecoder(f).Decode(p); err != nil {
		return nil, e.Wrap(op, err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
