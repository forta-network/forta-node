package botreq

import (
	"github.com/forta-network/forta-core-go/protocol"
)

// TxRequest contains the request data.
type TxRequest struct {
	Original *protocol.EvaluateTxRequest
}

// BlockRequest contains the request data.
type BlockRequest struct {
	Original *protocol.EvaluateBlockRequest
}

// CombinationRequest contains the request data.
type CombinationRequest struct {
	Original *protocol.EvaluateAlertRequest
}
