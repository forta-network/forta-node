package agentgrpc

import (
	"errors"

	"github.com/forta-network/forta-core-go/protocol"
)

// Error makes single error from our common errors defined in protobuf.
func Error(respErrs []*protocol.Error) error {
	var errMsg string
	for i, respErr := range respErrs {
		if i > 0 {
			errMsg += ", "
		}
		errMsg += respErr.Message
	}
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}
