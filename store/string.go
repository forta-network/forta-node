package store

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type StringStore interface {
	Get() (string, error)
	Put(string) error
}

type fileStringStore struct {
	path string
}

func NewFileStringStore(path string) *fileStringStore {
	return &fileStringStore{path: path}
}

func (fss *fileStringStore) Put(body string) error {
	return ioutil.WriteFile(fss.path, []byte(body), 0644)
}

func (fss *fileStringStore) Get() (string, error) {
	b, err := ioutil.ReadFile(fss.path)
	if err != nil {
		log.WithError(err).Warn("failed to read the last batch file")
		return "", nil
	}
	return strings.TrimSpace(string(b)), nil
}
