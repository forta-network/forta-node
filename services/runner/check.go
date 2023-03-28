package runner

import (
	"fmt"
	"net/url"
)

// CheckProxyAgainstScan checks given proxy URL against the scan API.
// The proxy API must specify as HTTP(s) when the scan API is WebSocket.
func CheckProxyAgainstScan(scan, proxy string) error {
	scanUrl, err := url.Parse(scan)
	if err != nil {
		return fmt.Errorf("invalid scan api url: %v", err)
	}
	// nothing to check if scan is already not websocket
	if !(scanUrl.Scheme == "ws" || scanUrl.Scheme == "wss") {
		return nil
	}

	if len(proxy) == 0 {
		return ErrBadProxyAPI
	}
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return fmt.Errorf("invalid proxy api url: %v", err)
	}
	// and proxy api must be either http or https
	if !(proxyUrl.Scheme == "http" || proxyUrl.Scheme == "https") {
		return ErrBadProxyAPI
	}
	return nil
}
