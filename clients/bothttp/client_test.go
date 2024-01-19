package bothttp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	r := require.New(t)

	mux := http.NewServeMux()

	// initially empty - no errors
	var respData HealthResponse

	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(respData)
		w.Write(b)
	}))
	server := http.Server{
		Addr:    "localhost:8183",
		Handler: mux,
	}
	go server.ListenAndServe()

	client := NewClient("localhost", 8183)
	_, err := client.Health(context.Background())
	r.NoError(err)

	respData.Errors = append(respData.Errors, "some error msg")
	_, err = client.Health(context.Background())
	r.Error(err)

	respData = HealthResponse{
		Metrics: []Metrics{
			{
				ChainID: 1,
				DataPoints: map[string][]float64{
					"tx.success": {1, 2, 3},
				},
			},
			{
				ChainID: 2,
				DataPoints: map[string][]float64{
					"tx.success": {3},
				},
			},
		},
	}

	hook := test.NewGlobal()

	metrics, err := client.Health(context.Background())
	r.NoError(err)
	r.EqualValues(respData.Metrics, metrics)

	responseSizeLimit = 1

	_, err = client.Health(context.Background())
	r.NoError(err)
	r.Equal(1, len(hook.Entries))
	r.Equal("response size limit for health check is reached", hook.LastEntry().Message)

	server.Close()
}
