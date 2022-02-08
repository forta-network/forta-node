package store

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ipfs/go-cid"
	log "github.com/sirupsen/logrus"
)

const lastBatchFileName = ".last-batch"

// BatchRefStore writes to and reads from somewhere the last batch reference.
type BatchRefStore interface {
	GetLast() (string, error)
	Put(string) error
}

type batchRefStore struct {
	filePath string
}

// NewBatchRefStore creates a new ref store.
func NewBatchRefStore(dir string) *batchRefStore {
	return &batchRefStore{
		filePath: path.Join(dir, lastBatchFileName),
	}
}

func (store *batchRefStore) GetLast() (string, error) {
	b, err := ioutil.ReadFile(store.filePath)
	if err != nil {
		log.WithError(err).Warn("failed to read the last batch file")
		return "", nil
	}
	if _, err = cid.Parse(string(b)); err != nil {
		return "", fmt.Errorf("invalid batch ref found: %v", err)
	}
	return strings.TrimSpace(string(b)), nil
}

func (store *batchRefStore) Put(ref string) error {
	if _, err := cid.Parse(ref); err != nil {
		return fmt.Errorf("invalid batch ref provided: %v", err)
	}
	return ioutil.WriteFile(store.filePath, []byte(ref), 0644)
}
