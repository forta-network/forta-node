package utils

import "net/http"

// BasicAuthTransport wraps the default transport with basic auth.
type BasicAuthTransport struct {
	username string
	password string
}

// NewBasicAuthTransport creates a new basic auth transport.
func NewBasicAuthTransport(username, password string) http.RoundTripper {
	return &BasicAuthTransport{
		username: username,
		password: password,
	}
}

// RoundTrip implements http.RoundTripper.
func (bat *BasicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(bat.username, bat.password)
	return http.DefaultTransport.RoundTrip(req)
}
