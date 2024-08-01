package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/alexKudryavtsev-web/grace_links_tg_bot/lib/e"
	"github.com/alexKudryavtsev-web/grace_links_tg_bot/storage"
)

const (
	defaultPerm = 0774
)

var (
	ErrNoSavedPage = errors.New("no saved page")
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(filePath, fName))
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)

	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, ErrNoSavedPage
	}

	file := files[rand.Intn(len(files))]

	return decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := p.Hash()

	if err != nil {
		return e.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove page: %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := p.Hash()

	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil 
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)
	default:
		return true, nil
	}
}

// UTILS
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}
