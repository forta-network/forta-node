package store

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"OpenZeppelin/zephyr-node/protocol"
)

func TestBadgerAlertStore_AddAlert(t *testing.T) {
	store, err := NewBadgerAlertStore()
	defer func() { assert.NoError(t, store.Clear()) }()

	assert.NoError(t, err)

	alertTime := time.Now()
	a := &protocol.Alert{Id: "test", Timestamp: alertTime.String()}
	assert.NoError(t, store.AddAlert(a))

	startRange := alertTime.Add(-1 * time.Minute)
	endRange := alertTime.Add(1 * time.Minute)

	res, err := store.GetAlerts(startRange, endRange)
	assert.NoError(t, store.AddAlert(a))

	assert.Lenf(t, res, 1, "result should have one entry")
}
