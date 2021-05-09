package store

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"

	"OpenZeppelin/fortify-node/protocol"
)

//DBPath is a local location of badger db (/db is a mounted volume)
const DBPath = "/db/fortify-alerts"

type AlertStore interface {
	GetAlerts(from time.Time, to time.Time) ([]*protocol.Alert, error)
	AddAlert(a *protocol.Alert) error
}

type BadgerAlertStore struct {
	db *badger.DB
}

func alertKey(a *protocol.Alert) string {
	return fmt.Sprintf("%s-%s", a.Timestamp, a.Id)
}

func prefixKey(t time.Time) string {
	return fmt.Sprintf("%s", t)
}

func isBetween(key []byte, startKey []byte, endKey []byte) bool {
	return string(key) >= string(startKey) && string(key) < string(endKey)
}

func (s *BadgerAlertStore) GetAlerts(from time.Time, to time.Time) ([]*protocol.Alert, error) {
	var result []*protocol.Alert
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		startKey := []byte(prefixKey(from))
		endKey := []byte(prefixKey(to))
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if !isBetween(item.Key(), startKey, endKey) {
				return nil
			}
			err := item.Value(func(v []byte) error {
				var alert protocol.Alert
				if err := proto.Unmarshal(v, &alert); err != nil {
					return err
				}
				result = append(result, &alert)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *BadgerAlertStore) Clear() error {
	return s.db.DropAll()
}

func (s *BadgerAlertStore) AddAlert(a *protocol.Alert) error {
	return s.db.Update(func(txn *badger.Txn) error {
		b, err := proto.Marshal(a)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(alertKey(a)), b).WithTTL(time.Hour * 24 * 7)
		err = txn.SetEntry(e)
		if err != nil {
			return err
		}
		return txn.SetEntry(e)
	})
}

func NewBadgerAlertStore() (*BadgerAlertStore, error) {
	db, err := badger.Open(badger.DefaultOptions("/tmp/fortify-alerts"))
	if err != nil {
		return nil, err
	}
	return &BadgerAlertStore{db: db}, nil
}
