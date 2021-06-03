package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"OpenZeppelin/fotify-node/protocol"
)

func TestFormat(t *testing.T) {
	val := "2021-05-11T23:41:11.304Z"
	time.Parse(AlertTimeFormat, val)
}

func TestBadgerAlertStore_AddAlert(t *testing.T) {
	store, err := NewBadgerAlertStoreWithPath("/tmp/alert-test")
	defer func() { assert.NoError(t, store.Clear()) }()

	assert.NoError(t, err)

	alertTime := time.Now().UTC()
	a := &protocol.SignedAlert{Alert: &protocol.Alert{Id: "test1", Timestamp: alertTime.Format(AlertTimeFormat)}}
	b := &protocol.SignedAlert{Alert: &protocol.Alert{Id: "test2", Timestamp: alertTime.Format(AlertTimeFormat)}}

	assert.NoError(t, store.AddAlert(a))
	assert.NoError(t, store.AddAlert(b))

	startRange := alertTime.Add(-1 * time.Minute)
	endRange := alertTime.Add(1 * time.Minute)

	res, err := store.QueryAlerts(&AlertQueryRequest{
		StartTime: startRange,
		EndTime:   endRange,
		PageToken: "",
		Limit:     1,
	})
	assert.NoError(t, err)

	expectedKeyB, err := alertKey(b)
	assert.NoError(t, err)
	assert.Lenf(t, res.Alerts, 1, "result should have one entry")
	assert.Equal(t, expectedKeyB, res.NextPageToken)

	res2, err := store.QueryAlerts(&AlertQueryRequest{
		StartTime: startRange,
		EndTime:   endRange,
		PageToken: res.NextPageToken,
		Limit:     1,
	})
	assert.NoError(t, err)
	assert.Lenf(t, res2.Alerts, 1, "result should have one entry")
}
