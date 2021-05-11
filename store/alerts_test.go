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

	alertTime := time.Now()
	a := &protocol.Alert{Id: "test1", Timestamp: alertTime.String()}
	b := &protocol.Alert{Id: "test2", Timestamp: alertTime.String()}

	assert.NoError(t, store.AddAlert(a))
	assert.NoError(t, store.AddAlert(b))

	startRange := alertTime.Add(-1 * time.Minute)
	endRange := alertTime.Add(1 * time.Minute)

	res, err := store.GetAlerts(AlertQueryRequest{
		FromTime:  startRange,
		ToTime:    endRange,
		PageStart: "",
		Limit:     1,
	})
	assert.NoError(t, err)
	assert.Lenf(t, res.Alerts, 1, "result should have one entry")
	assert.Equal(t, fmt.Sprintf("%s-%s", alertTime.String(), b.Id), res.NextPageStart)
}
