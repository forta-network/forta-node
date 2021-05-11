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
	a := &protocol.Alert{Id: "test1", Timestamp: alertTime.Format(time.RFC3339)}
	b := &protocol.Alert{Id: "test2", Timestamp: alertTime.Format(time.RFC3339)}

	assert.NoError(t, store.AddAlert(a))
	assert.NoError(t, store.AddAlert(b))

	startRange := alertTime.Add(-1 * time.Minute)
	endRange := alertTime.Add(1 * time.Minute)

	res, err := store.QueryAlerts(&AlertQueryRequest{
		FromTime:  startRange,
		ToTime:    endRange,
		PageToken: "",
		Limit:     1,
	})
	assert.NoError(t, err)
	assert.Lenf(t, res.Alerts, 1, "result should have one entry")
	assert.Equal(t, fmt.Sprintf("%s-%s", alertTime.Format(time.RFC3339), b.Id), res.NextPageToken)
}
