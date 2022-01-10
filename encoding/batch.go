package encoding

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/forta-protocol/forta-node/protocol"
	"github.com/golang/protobuf/proto"
)

func gzipBytes(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(b)
	if err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gunzipBytes(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func DecodeBatch(encoded string) (*protocol.AlertBatch, error) {
	zipped, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to base64decode batch: %v", err)
	}
	b, err := gunzipBytes(zipped)
	if err != nil {
		return nil, err
	}
	var batch protocol.AlertBatch
	if err := proto.Unmarshal(b, &batch); err != nil {
		return nil, err
	}
	return &batch, nil
}

func EncodeBatch(batch *protocol.AlertBatch) (string, error) {
	b, err := proto.Marshal(batch)
	if err != nil {
		return "", fmt.Errorf("failed to marshal batch: %v", err)
	}
	zipped, err := gzipBytes(b)
	if err != nil {
		return "", fmt.Errorf("failed to gzip batch: %v", err)
	}

	return base64.StdEncoding.EncodeToString(zipped), nil
}
