package agentgrpc

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

var (
	defaultCodec    = encoding.GetCodec(proto.Name)
	destPrepMsgType = reflect.TypeOf(&grpc.PreparedMsg{})
)

type preparedMsg struct {
	encodedData []byte
	hdr         []byte
	payload     []byte
}

// EncodeMessage encodes request as a PreparedMsg so the client stream can use it
// directly instead of allocating a new encoded message.
//
// See https://github.com/grpc/grpc-go/blob/1ffd63de37de4571028efedb6422e29d08716d0c/stream.go#L1623
func EncodeMessage(msg interface{}) (*grpc.PreparedMsg, error) {
	msgB, err := defaultCodec.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("agentgrpc: failed to encode message: %v", err)
	}
	hdr := make([]byte, 5)
	// write length of payload into header buffer
	binary.BigEndian.PutUint32(hdr[1:], uint32(len(msgB)))
	// hacky conversion to avoid compiler error
	return (*grpc.PreparedMsg)((unsafe.Pointer)(&preparedMsg{
		encodedData: msgB,
		payload:     msgB,
		hdr:         hdr,
	})), nil
}
