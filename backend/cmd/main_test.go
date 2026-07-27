package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestMain_ServesHTTP is a smoke test for the compiled binary: it builds
// cmd/main.go, runs it as a real subprocess against a live database, and
// checks that it actually starts listening and serves HTTP. This is the
// only test that exercises the process as a whole (env parsing, DB
// connection, service wiring, HTTP server) rather than individual pieces.
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

	resp := waitForServer(t, "http://127.0.0.1:"+port+"/", 5*time.Second)
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if got, want := string(body), "ok\n"; got != want {
		t.Errorf("GET / body = %q, want %q", got, want)
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
