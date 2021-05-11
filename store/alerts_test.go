package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"OpenZeppelin/fortify-node/protocol"
)

func TestBadgerAlertStore_AddAlert(t *testing.T) {
	store, err := NewBadgerAlertStoreWithPath("/tmp/alert-test")
	defer func() { assert.NoError(t, store.Clear()) }()

	assert.NoError(t, err)

	alertTime := time.Now().UTC()
	a := &protocol.Alert{Id: "test1", Timestamp: alertTime.Format(AlertTimeFormat)}
	b := &protocol.Alert{Id: "test2", Timestamp: alertTime.Format(AlertTimeFormat)}

	assert.NoError(t, store.AddAlert(a))
	assert.NoError(t, store.AddAlert(b))

	startRange := alertTime.Add(-1 * time.Minute)
	endRange := alertTime.Add(1 * time.Minute)

	ks, err := store.GetAllKeys()
	assert.NoError(t, err)
	for _, k := range ks {
		t.Log(k)
	}

	res, err := store.QueryAlerts(&AlertQueryRequest{
		StartTime: startRange,
		EndTime:   endRange,
		PageToken: "",
		Limit:     1,
	})
	assert.NoError(t, err)
	assert.Lenf(t, res.Alerts, 1, "result should have one entry")
	assert.Equal(t, fmt.Sprintf("%s-%s", alertTime.Format(AlertTimeKeyFormat), b.Id), res.NextPageToken)

	res2, err := store.QueryAlerts(&AlertQueryRequest{
		StartTime: startRange,
		EndTime:   endRange,
		PageToken: res.NextPageToken,
		Limit:     1,
	})
	assert.NoError(t, err)
	assert.Lenf(t, res2.Alerts, 1, "result should have one entry")
}
