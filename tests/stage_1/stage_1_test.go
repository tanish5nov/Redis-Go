package stage_1_test

import (
	"net"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func buildServerBinary(t *testing.T) string {
	t.Helper()

	rootDir, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to resolve root dir: %v", err)
	}

	binPath := filepath.Join(t.TempDir(), "redis-server")
	buildCmd := exec.Command("go", "build", "-o", binPath, ".")
	buildCmd.Dir = rootDir
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server binary: %v\nOutput:\n%s", err, string(out))
	}

	return binPath
}

func startServer(t *testing.T) (*exec.Cmd, func()) {
	t.Helper()

	binPath := buildServerBinary(t)
	cmd := exec.Command(binPath)

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	cleanup := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 150*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return cmd, cleanup
		}
		time.Sleep(50 * time.Millisecond)
	}

	cleanup()
	t.Fatalf("Server did not bind and listen on 127.0.0.1:6379 within 3 seconds")
	return nil, nil
}

func TestServerBindsAndAcceptsConnection(t *testing.T) {
	_, cleanup := startServer(t)
	defer cleanup()

	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
	if err != nil {
		t.Fatalf("Expected connection to 127.0.0.1:6379 to succeed, got: %v", err)
	}
	_ = conn.Close()
}

func TestServerStaysAliveForMultipleSequentialConnections(t *testing.T) {
	_, cleanup := startServer(t)
	defer cleanup()

	for i := 0; i < 5; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
		if err != nil {
			t.Fatalf("Sequential client %d failed to connect: %v", i+1, err)
		}
		_ = conn.Close()
		time.Sleep(30 * time.Millisecond)
	}
}

func TestServerAcceptsRapidConcurrentConnections(t *testing.T) {
	_, cleanup := startServer(t)
	defer cleanup()

	const numClients = 20
	var wg sync.WaitGroup
	errs := make(chan error, numClients)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 2*time.Second)
			if err != nil {
				errs <- err
				return
			}
			_ = conn.Close()
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("Concurrent client connection failed: %v", err)
	}
}

func TestServerHandlesPortCollisionWithoutPanic(t *testing.T) {
	dummyListener, err := net.Listen("tcp", ":6379")
	if err != nil {
		t.Fatalf("Failed to bind dummy listener to :6379: %v", err)
	}
	defer dummyListener.Close()

	binPath := buildServerBinary(t)
	cmd := exec.Command(binPath)

	output, _ := cmd.CombinedOutput()
	outStr := string(output)

	for _, badSignal := range []string{"panic:", "nil pointer dereference", "SIGSEGV"} {
		if strings.Contains(outStr, badSignal) {
			t.Fatalf("Server panicked on startup error:\n%s", outStr)
		}
	}
}
