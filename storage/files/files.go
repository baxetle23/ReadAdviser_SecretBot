package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"telegram_bot/lib"
	"telegram_bot/storage"
	"time"
)

type StorageFiles struct {
	basePath string
}

const (
	defaultPerm = 0774
)

func New(basePath string) *StorageFiles {
	return &StorageFiles{
		basePath: basePath,
	}
}

func (s *StorageFiles) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = lib.WrapIfError("can't pick random reference", err) }()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s *StorageFiles) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return lib.WrapOnError("Can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return lib.WrapOnError("can't remove file", err)
	}

	return nil
}

func (s *StorageFiles) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, lib.WrapOnError("can't open file", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, lib.WrapOnError("can't decode page", err)
	}
	return &p, nil
}

func (s *StorageFiles) Save(page *storage.Page) (err error) {
	defer func() { err = lib.WrapIfError("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
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

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s *StorageFiles) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, lib.WrapOnError("Can't remove file", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, lib.WrapOnError(msg, err)
	}
	return true, nil
}
