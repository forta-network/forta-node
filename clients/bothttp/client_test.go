package bothttp

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

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

	// TODO: check returned metrics

	server.Close()
}
