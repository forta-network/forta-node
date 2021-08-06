package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/golang/protobuf/proto"

	"OpenZeppelin/fortify-node/protocol"
)

//DBPath is a local location of badger db (/db is a mounted volume)
const DBPath = "/db/fortify-alerts"
const AlertTimeFormat = time.RFC3339Nano

type Operator string

const Equals Operator = "EQUALS"
const NotEquals Operator = "NOT_EQUALS"

var ErrNoPruneNeeded = errors.New("no prune was deemed necessary")

type FilterCriterion struct {
	Operator Operator
	Field    string
	Values   []string
}

func stringVal(fieldName string, alert *protocol.Alert) (string, bool) {
	if fieldName == "alertId" {
		return alert.Finding.AlertId, true
	}
	if fieldName == "severity" {
		return alert.Finding.Severity.String(), true
	}
	if fieldName == "protocol" {
		return alert.Finding.Protocol, true
	}
	if fieldName == "everestId" {
		return alert.Finding.EverestId, true
	}
	if fieldName == "type" {
		return alert.Finding.Type.String(), true
	}
	if fieldName == "agentName" {
		return alert.Agent.Name, true
	}
	if fieldName == "agentImage" {
		return alert.Agent.Image, true
	}
	if fieldName == "agentImageHash" {
		return alert.Agent.ImageHash, true
	}
	if val, ok := alert.Tags[fieldName]; ok {
		return val, true
	}
	if val, ok := alert.Finding.Metadata[fieldName]; ok {
		return val, true
	}
	return "", false
}

func (fc *FilterCriterion) Matches(alert *protocol.SignedAlert) bool {
	val, ok := stringVal(fc.Field, alert.Alert)
	if ok && fc.Operator == Equals {
		// any must match
		for _, v := range fc.Values {
			if v == val {
				return true
			}
		}
	} else if fc.Operator == NotEquals {
		if !ok {
			return true
		}
		// all must not match
		for _, v := range fc.Values {
			if v == val {
				return false
			}
		}
		return true
	}
	return false
}

type AlertQueryRequest struct {
	StartTime time.Time
	EndTime   time.Time
	PageToken string
	Limit     int
	Reverse   bool
	Criteria  []*FilterCriterion
}

func (r *AlertQueryRequest) Matches(alert *protocol.SignedAlert) bool {
	// must match all, if any
	for _, fc := range r.Criteria {
		if !fc.Matches(alert) {
			return false
		}
	}
	return true
}

func (r *AlertQueryRequest) Json() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

type AlertQueryResponse struct {
	Alerts        []*protocol.SignedAlert
	NextPageToken string
}

type AlertStore interface {
	QueryAlerts(request *AlertQueryRequest) (*AlertQueryResponse, error)
	AddAlert(a *protocol.SignedAlert) error
	Prune() error
}

type BadgerAlertStore struct {
	db *badger.DB
}

func alertKey(a *protocol.SignedAlert) (string, error) {
	ts, err := time.Parse(AlertTimeFormat, a.Alert.Timestamp)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", formatSearchKey(ts), a.Alert.Id), nil
}

func formatSearchKey(t time.Time) string {
	return fmt.Sprintf("%d", t.UnixNano()/1e6)
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

func (s *BadgerAlertStore) Prune() error {
	err := s.db.RunValueLogGC(0.5)
	if err == badger.ErrNoRewrite {
		return ErrNoPruneNeeded
	}
	return err
}

func (s *BadgerAlertStore) QueryAlerts(request *AlertQueryRequest) (*AlertQueryResponse, error) {
	result := &AlertQueryResponse{
		Alerts: make([]*protocol.SignedAlert, 0),
	}

	// seek to this key first
	startKey := request.PageToken
	if startKey == "" {
		if request.Reverse {
			startKey = formatSearchKey(request.EndTime)
		} else {
			startKey = formatSearchKey(request.StartTime)
		}
	}

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		opts.Reverse = request.Reverse
		it := txn.NewIterator(opts)
		defer it.Close()

		startTime := []byte(formatSearchKey(request.StartTime))
		endTime := []byte(formatSearchKey(request.EndTime))

		for it.Seek([]byte(startKey)); it.Valid(); it.Next() {
			item := it.Item()
			if len(result.Alerts) == request.Limit {
				result.NextPageToken = string(item.Key())
				return nil
			}
			if !isBetween(item.Key(), startTime, endTime) {
				return nil
			}
			err := item.Value(func(v []byte) error {
				var alert protocol.SignedAlert
				if err := proto.Unmarshal(v, &alert); err != nil {
					return err
				}
				if request.Matches(&alert) {
					result.Alerts = append(result.Alerts, &alert)
				}
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

func (s *BadgerAlertStore) AddAlert(a *protocol.SignedAlert) error {
	return s.db.Update(func(txn *badger.Txn) error {
		b, err := proto.Marshal(a)
		if err != nil {
			return err
		}
		ak, err := alertKey(a)
		if err != nil {
			return err
		}
		e := badger.NewEntry([]byte(ak), b).WithTTL(time.Hour * 24 * 7)
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
