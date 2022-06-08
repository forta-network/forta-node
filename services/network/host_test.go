package network

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadHostNetworking(t *testing.T) {
	r := require.New(t)

	host := &Host{
		DefaultInterfaceName: "eth0",
		DefaultSubnet:        "192.168.0.0/24",
		DefaultGateway:       "192.168.0.1",
		Docker0Subnet:        "10.99.0.0/24",
	}

	var buf bytes.Buffer
	WriteHostNetworking(&buf, host)

	detected := UnmarshalHostNetworking(string(buf.Bytes()))
	r.Equal(host.DefaultInterfaceName, detected.DefaultInterfaceName)
	r.Equal(host.DefaultSubnet, detected.DefaultSubnet)
	r.Equal(host.DefaultGateway, detected.DefaultGateway)
	r.Equal(host.Docker0Subnet, detected.Docker0Subnet)
}
