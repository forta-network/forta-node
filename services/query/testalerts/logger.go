package testalerts

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/forta-protocol/forta-node/protocol"
)

// Logger logs the test alerts either to a log file or to a webhook.
type Logger struct {
	file       *os.File
	webhookUrl string
}

// NewLogger creates a new logger.
func NewLogger(dst string) *Logger {
	if len(dst) != 0 {
		return newWebhookLogger(dst)
	}

	if err := os.MkdirAll("/test-alerts", 0666); err != nil {
		panic(fmt.Errorf("failed to create the test alerts dir: %v", err))
	}
	file, err := os.Create(fmt.Sprintf("/test-alerts/forta-test-alert-log-%d", time.Now().Unix()))
	if err != nil {
		panic(fmt.Errorf("failed to create the test alert log file: %v", err))
	}
	if err != nil {
		panic(err)
	}
	return &Logger{file: file}
}

func newWebhookLogger(dst string) *Logger {
	u, err := url.Parse(dst)
	if err != nil {
		panic(fmt.Errorf("failed to parse the webhook url: %v", err))
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		panic("non-http webhook url")
	}
	return &Logger{webhookUrl: dst}
}

// Close implemenets io.Closer.
func (logger *Logger) Close() error {
	if logger.file == nil {
		return nil
	}
	return logger.file.Close()
}

// LogTestAlert logs the test alert by marshalling to JSON.
func (logger *Logger) LogTestAlert(ctx context.Context, alert *protocol.SignedAlert) error {
	b, _ := json.Marshal(alert)
	if logger.file != nil {
		_, err := fmt.Fprintln(logger.file, string(b))
		if err != nil {
			return fmt.Errorf("failed to write to the test alert file: %v", err)
		}
	}
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "POST", logger.webhookUrl, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("failed to send the test alert: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("test alert webhook request failed: %v", err)
	}
	return nil
}
