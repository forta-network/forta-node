package runner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckProxyAgainstScan(t *testing.T) {
	testCases := []struct {
		name  string
		scan  string
		proxy string
		valid bool
	}{
		{
			name:  "scan is http",
			scan:  "http://foo.bar",
			proxy: "",
			valid: true,
		},
		{
			name:  "scan is websocket, proxy is empty",
			scan:  "wss://foo.bar",
			proxy: "",
			valid: false,
		},
		{
			name:  "scan is websocket, proxy is websocket",
			scan:  "wss://foo.bar",
			proxy: "wss://foo.bar",
			valid: false,
		},
		{
			name:  "scan is websocket, proxy is http",
			scan:  "wss://foo.bar",
			proxy: "http://any.thing",
			valid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			r := require.New(t)

			err := CheckProxyAgainstScan(testCase.scan, testCase.proxy)
			if testCase.valid {
				r.NoError(err)
			} else {
				r.Error(err)
			}
		})
	}
}
