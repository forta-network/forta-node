package cooldown

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testID = "1"

func TestCooldown(t *testing.T) {
	r := require.New(t)

	coolDown := &cooldown{
		threshold:        1,
		cooldownDuration: time.Hour,
		counters:         make(map[string]*cooldownCounter),
	}
	r.False(coolDown.ShouldCoolDown(testID))
	r.True(coolDown.ShouldCoolDown(testID))
	coolDown.counters[testID].cooldownEndsAt = time.Now().Add(-time.Hour)
	r.False(coolDown.ShouldCoolDown(testID))
	r.True(coolDown.ShouldCoolDown(testID))
	r.True(coolDown.ShouldCoolDown(testID))

	coolDown.doCleanup()
	r.Len(coolDown.counters, 1)
}
