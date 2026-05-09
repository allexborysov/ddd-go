package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

type Env struct {
	BaseURL string
	Client  *http.Client
}

func NewEnv() *Env {
	base := os.Getenv("E2E_BASE_URL")
	if base == "" {
		base = "http://127.0.0.1:8080"
	}
	return &Env{
		BaseURL: base,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

var uniqueCounter uint64

// UniqueSuffix returns a short identifier safe for parallel use:
// millisecond timestamp + atomic counter, so two callers in the same
// millisecond still get distinct values. Length stays within MSN's
// 20-char limit (MSN- prefix + 13 millis + up to 3 digits = ≤20).
func UniqueSuffix() string {
	n := atomic.AddUint64(&uniqueCounter, 1) % 1000
	return fmt.Sprintf("%d%03d", time.Now().UnixMilli(), n)
}

// postJSON POSTs body to env.BaseURL+path and unmarshals a 2xx response
// into out. Any non-2xx status t.Fatals the calling test.
func postJSON(ctx context.Context, t *testing.T, env *Env, path string, body any, out any) {
	t.Helper()

	buf, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	url := env.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := env.Client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	if resp.StatusCode >= 300 {
		t.Fatalf("POST %s: status %d, body=%s", url, resp.StatusCode, string(respBody))
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			t.Fatalf("unmarshal response (%s): %v\nbody=%s", url, err, string(respBody))
		}
	}
}

// postJSONExpectError asserts the server returned an error status (>= 400)
// and returns the parsed {"message": ...} body.
func postJSONExpectError(ctx context.Context, t *testing.T, env *Env, path string, body any) (int, string) {
	t.Helper()

	buf, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	url := env.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := env.Client.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}

	if resp.StatusCode < 400 {
		t.Fatalf("POST %s: expected error status, got %d, body=%s", url, resp.StatusCode, string(respBody))
	}

	var msg struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &msg); err != nil || msg.Message == "" {
		return resp.StatusCode, string(respBody)
	}
	return resp.StatusCode, msg.Message
}
