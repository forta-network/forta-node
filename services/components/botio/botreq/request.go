package botreq

import (
	"github.com/forta-network/forta-core-go/protocol"
	"google.golang.org/grpc"
)

// TxRequest contains the original request data and the encoded message.
type TxRequest struct {
	Original *protocol.EvaluateTxRequest
	Encoded  *grpc.PreparedMsg
}

// BlockRequest contains the original request data and the encoded message.
type BlockRequest struct {
	Original *protocol.EvaluateBlockRequest
	Encoded  *grpc.PreparedMsg
}

// CombinationRequest contains the original request data and the encoded message.
type CombinationRequest struct {
	Original *protocol.EvaluateAlertRequest
	Encoded  *grpc.PreparedMsg
}
