package healthutils

import log "github.com/sirupsen/logrus"

// DefaultHealthServerErrHandler handlers health server error.
func DefaultHealthServerErrHandler(err error) {
	log.WithError(err).Panic("health server failed")
}
