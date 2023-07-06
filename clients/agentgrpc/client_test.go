package agentgrpc

import (
	"errors"
	"testing"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func Test_evaluateHealthCheckResult(t *testing.T) {
	// Test case 1: No errors
	resp := &protocol.HealthCheckResponse{}
	err := evaluateHealthCheckResult(nil, resp)
	assert.Nil(t, err, "Expected no error")

	// Test case 2: Invocation error
	invokeErr := errors.New("invoke error")
	expectedError := multierror.Append(invokeErr).Error()
	err = evaluateHealthCheckResult(invokeErr, new(protocol.HealthCheckResponse))
	assert.EqualError(t, err, expectedError, "Unexpected error")

	// Test case 3: Response errors
	resp = &protocol.HealthCheckResponse{
		Errors: []*protocol.Error{
			{Message: "error 1"},
			{Message: "error 2"},
		},
	}
	err = evaluateHealthCheckResult(nil, resp)
	expectedError = "2 errors occurred:\n\t* error 1\n\t* error 2\n\n"
	assert.EqualError(t, err, expectedError, "Unexpected error")

	// Test case 4: Invocation and response errors
	err = evaluateHealthCheckResult(invokeErr, resp)
	expectedError = "3 errors occurred:\n\t* invoke error\n\t* error 1\n\t* error 2\n\n"
	assert.EqualError(t, err, expectedError, "Unexpected error")
}
