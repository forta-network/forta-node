package agentgrpc

import (
	"errors"
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsHealthCheckSuccess(t *testing.T) {
	// Test case 1: Successful health check
	resp := &protocol.HealthCheckResponse{
		Status: protocol.HealthCheckResponse_SUCCESS,
	}
	success := isHealthCheckSuccess(nil, resp)
	assert.True(t, success, "Expected health check to be successful")

	// Test case 2: Unimplemented method error
	invokeErr := status.New(codes.Unimplemented, "").Err()
	success = isHealthCheckSuccess(invokeErr, new(protocol.HealthCheckResponse))
	assert.True(t, success, "Expected health check to be successful (unimplemented)")

	// Test case 3: Other invocation error
	invokeErr = errors.New("other error")
	success = isHealthCheckSuccess(invokeErr, nil)
	assert.False(t, success, "Expected health check to be unsuccessful")
}

func TestExtractHealthCheckError(t *testing.T) {
	// Test case 1: No errors
	resp := &protocol.HealthCheckResponse{}
	err := extractHealthCheckError(nil, resp)
	assert.Nil(t, err, "Expected no error")

	// Test case 2: Invocation error
	invokeErr := errors.New("invoke error")
	expectedError := multierror.Append(invokeErr).Error()
	err = extractHealthCheckError(invokeErr, new(protocol.HealthCheckResponse))
	assert.EqualError(t, err, expectedError, "Unexpected error")

	// Test case 3: Response errors
	resp = &protocol.HealthCheckResponse{
		Errors: []*protocol.Error{
			{Message: "error 1"},
			{Message: "error 2"},
		},
	}
	err = extractHealthCheckError(nil, resp)
	expectedError = "2 errors occurred:\n\t* error 1\n\t* error 2\n\n"
	t.Log(err)
	assert.EqualError(t, err, expectedError, "Unexpected error")

	// Test case 4: Invocation and response errors
	err = extractHealthCheckError(invokeErr, resp)
	expectedError = "3 errors occurred:\n\t* invoke error\n\t* error 1\n\t* error 2\n\n"
	assert.EqualError(t, err, expectedError, "Unexpected error")
}
