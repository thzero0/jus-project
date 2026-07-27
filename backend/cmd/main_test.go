package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

// TestMain_ServesHTTP is a smoke test for the compiled binary: it builds
// cmd/main.go, runs it as a real subprocess against a live database, and
// sends a real GraphQL query over HTTP. This is the only test that
// exercises the process as a whole — env parsing, DB connection, service
// wiring, and the GraphQL layer together — rather than individual pieces.
//
// It requires DATABASE_URL to point at a running, seeded database — see
// docs/testing-backend.md. It is skipped when DATABASE_URL is unset, so
// `go test ./...` stays green in CI, which doesn't provision a database.
func TestMain_ServesHTTP(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set; see docs/testing-backend.md to run this test")
	}

	port := freePort(t)

	binPath := filepath.Join(t.TempDir(), "api")
	build := exec.Command("go", "build", "-o", binPath, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath)
	cmd.Env = append(os.Environ(), "DATABASE_URL="+databaseURL, "PORT="+port)
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting binary: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	baseURL := "http://127.0.0.1:" + port
	readyResp := waitForServer(t, baseURL+"/", 5*time.Second)
	_ = readyResp.Body.Close()

	// "the elder scrolls vi" (not just "the e") deliberately narrows to a
	// single, known match: dozens of other "The Elder Scrolls ..." titles
	// exist in the dataset, and alphabetical sort (':' sorts before 'I')
	// pushes "The Elder Scrolls VI" past the 20-result cap for a shorter
	// prefix like "the e".
	const term = "the elder scrolls vi"
	reqBody, err := json.Marshal(map[string]any{
		"query":     `query($term: String!) { suggestions(term: $term) }`,
		"variables": map[string]any{"term": term},
	})
	if err != nil {
		t.Fatalf("marshaling GraphQL request: %v", err)
	}

	resp, err := http.Post(baseURL+"/graphql", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("POST /graphql: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("POST /graphql status = %d, body = %s", resp.StatusCode, body)
	}

	var result struct {
		Data struct {
			Suggestions []string `json:"suggestions"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decoding GraphQL response: %v", err)
	}
	if len(result.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", result.Errors)
	}

	want := []string{"The Elder Scrolls VI"}
	if !slices.Equal(result.Data.Suggestions, want) {
		t.Errorf("suggestions(term: %q) = %v, want %v", term, result.Data.Suggestions, want)
	}
}

func freePort(t *testing.T) string {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding a free port: %v", err)
	}
	defer func() { _ = l.Close() }()

	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		t.Fatalf("parsing free port: %v", err)
	}
	return port
}

func waitForServer(t *testing.T, url string, timeout time.Duration) *http.Response {
	t.Helper()

	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			return resp
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server at %s did not become ready: %v", url, lastErr)
	return nil
}
