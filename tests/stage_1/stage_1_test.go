package stage_1_test

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func startServer(t *testing.T) (*exec.Cmd, func()) {
	t.Helper()

	// Locate root dir relative to this test file
	rootDir, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to resolve root dir: %v", err)
	}

	binPath := filepath.Join(rootDir, "redis-server-test-bin")
	buildCmd := exec.Command("go", "build", "-o", binPath, "redisConnection.go")
	buildCmd.Dir = rootDir
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server binary: %v, output: %s", err, string(out))
	}

	cmd := exec.Command(binPath)
	cmd.Dir = rootDir

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	cleanup := func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
		_ = os.Remove(binPath)
	}

	// Try dialing port 6379 until server is listening or timeout expires
	var connected bool
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 150*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			connected = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !connected {
		cleanup()
		t.Fatalf("Server did not bind and listen on 127.0.0.1:6379 within 3 seconds")
	}

	return cmd, cleanup
}

func TestServerBindsAndAcceptsConnection(t *testing.T) {
	_, cleanup := startServer(t)
	defer cleanup()

	conn, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
	if err != nil {
		t.Fatalf("Expected connection to succeed on 127.0.0.1:6379, got error: %v", err)
	}
	defer conn.Close()
}

func TestServerStaysAliveForMultipleConnections(t *testing.T) {
	_, cleanup := startServer(t)
	defer cleanup()

	// Client 1 connects and disconnects
	conn1, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
	if err != nil {
		t.Fatalf("Client 1 failed to connect: %v", err)
	}
	_ = conn1.Close()

	// Short delay
	time.Sleep(100 * time.Millisecond)

	// Client 2 connects - server should still be accepting connections
	conn2, err := net.DialTimeout("tcp", "127.0.0.1:6379", 1*time.Second)
	if err != nil {
		t.Fatalf("Client 2 failed to connect (server may have exited after first connection): %v", err)
	}
	_ = conn2.Close()
}

func TestServerHandlesPortCollisionCleanly(t *testing.T) {
	// Occupy port 6379 with a dummy listener first
	dummyListener, err := net.Listen("tcp", ":6379")
	if err != nil {
		t.Fatalf("Failed to bind dummy listener to :6379: %v", err)
	}
	defer dummyListener.Close()

	rootDir, err := filepath.Abs("../../")
	if err != nil {
		t.Fatalf("Failed to resolve root dir: %v", err)
	}

	// Try running redisConnection.go when port is already occupied
	binPath := filepath.Join(rootDir, "redis-server-test-bin")
	cmd := exec.Command(binPath)
	cmd.Dir = rootDir

	output, err := cmd.CombinedOutput()
	// It should exit with an error, but it should NOT panic
	outStr := string(output)
	if outStr == "" && err == nil {
		t.Fatalf("Expected server to fail to start on an occupied port")
	}

	for _, badSignal := range []string{"panic:", "nil pointer dereference", "SIGSEGV"} {
		if bytesContains(outStr, badSignal) {
			t.Fatalf("Server crashed with unexpected panic on port collision: %s", outStr)
		}
	}
}

func bytesContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

