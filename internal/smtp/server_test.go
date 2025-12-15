package smtp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/lawnchairsociety/devsmtp/internal/config"
	"github.com/lawnchairsociety/devsmtp/internal/database"
)

func setupTestServer(t *testing.T) (*Server, *database.DB, *Logger, int, func()) {
	t.Helper()

	// Create temp database
	tmpFile, err := os.CreateTemp("", "devsmtp-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	db, err := database.New(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	// Find available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: port,
		},
	}

	logger := NewLogger(100)
	server := NewServer(cfg, db, logger)

	// Start server in background
	go func() {
		_ = server.ListenAndServe()
	}()

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)

	cleanup := func() {
		db.Close()
		os.Remove(tmpFile.Name())
	}

	return server, db, logger, port, cleanup
}

func connectToServer(t *testing.T, port int) net.Conn {
	t.Helper()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}

	return conn
}

func readLine(t *testing.T, conn net.Conn) string {
	t.Helper()
	if err := conn.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("failed to set read deadline: %v", err)
	}
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read line: %v", err)
	}
	return strings.TrimSpace(line)
}

func writeLine(t *testing.T, conn net.Conn, line string) {
	t.Helper()
	_, err := fmt.Fprintf(conn, "%s\r\n", line)
	if err != nil {
		t.Fatalf("failed to write line: %v", err)
	}
}

func TestServerGreeting(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	greeting := readLine(t, conn)
	if !strings.HasPrefix(greeting, "220") {
		t.Errorf("expected 220 greeting, got: %s", greeting)
	}
}

func TestHELO(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "HELO localhost")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response, got: %s", response)
	}
}

func TestEHLO(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "EHLO localhost")

	// Read all EHLO responses
	var responses []string
	reader := bufio.NewReader(conn)
	for {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		responses = append(responses, line)
		if strings.HasPrefix(line, "250 ") { // Last line has space, not dash
			break
		}
	}

	if len(responses) == 0 {
		t.Fatal("expected EHLO responses")
	}

	if !strings.HasPrefix(responses[0], "250") {
		t.Errorf("expected 250 response, got: %s", responses[0])
	}
}

func TestMAILFROM(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting
	writeLine(t, conn, "HELO localhost")
	readLine(t, conn)

	writeLine(t, conn, "MAIL FROM:<sender@example.com>")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response, got: %s", response)
	}
}

func TestRCPTTO(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting
	writeLine(t, conn, "HELO localhost")
	readLine(t, conn)
	writeLine(t, conn, "MAIL FROM:<sender@example.com>")
	readLine(t, conn)

	writeLine(t, conn, "RCPT TO:<recipient@example.com>")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response, got: %s", response)
	}
}

func TestRCPTTOWithoutMAIL(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting
	writeLine(t, conn, "HELO localhost")
	readLine(t, conn)

	writeLine(t, conn, "RCPT TO:<recipient@example.com>")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "503") {
		t.Errorf("expected 503 response, got: %s", response)
	}
}

func TestFullMessageFlow(t *testing.T) {
	_, db, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	readLineReader := func() string {
		_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		line, _ := reader.ReadString('\n')
		return strings.TrimSpace(line)
	}

	readLineReader() // greeting

	writeLine(t, conn, "HELO localhost")
	readLineReader()

	writeLine(t, conn, "MAIL FROM:<sender@test.com>")
	readLineReader()

	writeLine(t, conn, "RCPT TO:<recipient@test.com>")
	readLineReader()

	writeLine(t, conn, "DATA")
	response := readLineReader()
	if !strings.HasPrefix(response, "354") {
		t.Fatalf("expected 354 response, got: %s", response)
	}

	// Send message content
	writeLine(t, conn, "Subject: Test Email")
	writeLine(t, conn, "From: sender@test.com")
	writeLine(t, conn, "To: recipient@test.com")
	writeLine(t, conn, "")
	writeLine(t, conn, "This is the body of the email.")
	writeLine(t, conn, ".")

	response = readLineReader()
	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response after DATA, got: %s", response)
	}

	writeLine(t, conn, "QUIT")

	// Verify message was saved
	messages, err := db.GetMessages()
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	msg := messages[0]
	if msg.Subject != "Test Email" {
		t.Errorf("expected subject 'Test Email', got %q", msg.Subject)
	}
	if !strings.Contains(msg.Sender, "sender@test.com") {
		t.Errorf("expected sender to contain 'sender@test.com', got %q", msg.Sender)
	}
}

func TestNOOP(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "NOOP")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response, got: %s", response)
	}
}

func TestRSET(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting
	writeLine(t, conn, "HELO localhost")
	readLine(t, conn)
	writeLine(t, conn, "MAIL FROM:<sender@example.com>")
	readLine(t, conn)

	writeLine(t, conn, "RSET")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "250") {
		t.Errorf("expected 250 response, got: %s", response)
	}

	// After RSET, RCPT should fail (no MAIL FROM)
	writeLine(t, conn, "RCPT TO:<recipient@example.com>")
	response = readLine(t, conn)

	if !strings.HasPrefix(response, "503") {
		t.Errorf("expected 503 response after RSET, got: %s", response)
	}
}

func TestQUIT(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "QUIT")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "221") {
		t.Errorf("expected 221 response, got: %s", response)
	}
}

func TestVRFY(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "VRFY user@example.com")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "252") {
		t.Errorf("expected 252 response, got: %s", response)
	}
}

func TestUnknownCommand(t *testing.T) {
	_, _, _, port, cleanup := setupTestServer(t)
	defer cleanup()

	conn := connectToServer(t, port)
	defer conn.Close()

	readLine(t, conn) // greeting

	writeLine(t, conn, "INVALID")
	response := readLine(t, conn)

	if !strings.HasPrefix(response, "502") {
		t.Errorf("expected 502 response, got: %s", response)
	}
}

func TestLogger(t *testing.T) {
	logger := NewLogger(10)

	logger.Info("test message %d", 1)
	logger.Warn("warning message")
	logger.Error("error message")
	logger.Debug("debug message")

	// Read from channel
	entries := make([]LogEntry, 0)
	timeout := time.After(100 * time.Millisecond)

loop:
	for {
		select {
		case entry := <-logger.Channel():
			entries = append(entries, entry)
		case <-timeout:
			break loop
		}
	}

	if len(entries) != 4 {
		t.Errorf("expected 4 log entries, got %d", len(entries))
	}
}

func TestLogEntryString(t *testing.T) {
	entry := LogEntry{
		Time:    time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC),
		Level:   LogInfo,
		Message: "test message",
	}

	str := entry.String()
	if !strings.Contains(str, "10:30:45") {
		t.Errorf("expected time in string, got: %s", str)
	}
	if !strings.Contains(str, "INFO") {
		t.Errorf("expected INFO in string, got: %s", str)
	}
	if !strings.Contains(str, "test message") {
		t.Errorf("expected message in string, got: %s", str)
	}
}

func TestLogLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogInfo, "INFO"},
		{LogWarning, "WARN"},
		{LogError, "ERROR"},
		{LogDebug, "DEBUG"},
	}

	for _, tt := range tests {
		if tt.level.String() != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, tt.level.String())
		}
	}
}
