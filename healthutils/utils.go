package healthutils

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

// DefaultHealthServerErrHandler handlers health server error.
func DefaultHealthServerErrHandler(err error) {
	if strings.Contains(strings.ToLower(err.Error()), "server closed") {
		log.WithError(err).Warn("health server was shut down")
		return
	}
	log.WithError(err).Panic("health server failed")
}
