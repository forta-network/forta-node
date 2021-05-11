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

type AlertQueryRequest struct {
	FromTime  time.Time
	ToTime    time.Time
	PageToken string
	Limit     int
}

type AlertQueryResponse struct {
	Alerts        []*protocol.Alert
	NextPageToken string
}

type AlertStore interface {
	QueryAlerts(request *AlertQueryRequest) (*AlertQueryResponse, error)
	AddAlert(a *protocol.Alert) error
}

type BadgerAlertStore struct {
	db *badger.DB
}

func alertKey(a *protocol.Alert) string {
	return fmt.Sprintf("%s-%s", a.Timestamp, a.Id)
}

func prefixKey(t time.Time) string {
	return t.Format(time.RFC3339)
}

func isBetween(key []byte, startKey []byte, endKey []byte) bool {
	return string(key) >= string(startKey) && string(key) < string(endKey)
}

// GetAllKeys is a utility method for debugging
func (s *BadgerAlertStore) GetAllKeys() ([]string, error) {
	var keys []string
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, err
}

func (s *BadgerAlertStore) QueryAlerts(request *AlertQueryRequest) (*AlertQueryResponse, error) {
	result := &AlertQueryResponse{
		Alerts: make([]*protocol.Alert, 0),
	}
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		startKey := []byte(prefixKey(request.FromTime))
		if request.PageToken != "" {
			startKey = []byte(request.PageToken)
		}
		endKey := []byte(prefixKey(request.ToTime))
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if len(result.Alerts) == request.Limit {
				result.NextPageToken = string(item.Key())
				return nil
			}
			if !isBetween(item.Key(), startKey, endKey) {
				return nil
			}
			err := item.Value(func(v []byte) error {
				var alert protocol.Alert
				if err := proto.Unmarshal(v, &alert); err != nil {
					return err
				}
				result.Alerts = append(result.Alerts, &alert)
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
	db, err := badger.Open(badger.DefaultOptions("/db/fortify-alerts"))
	if err != nil {
		return nil, err
	}
	return &BadgerAlertStore{db: db}, nil
}

func NewBadgerAlertStoreWithPath(path string) (*BadgerAlertStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	return &BadgerAlertStore{db: db}, nil
}
