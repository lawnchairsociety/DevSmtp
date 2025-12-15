package database

import (
	"os"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*DB, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "devsmtp-test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	db, err := New(tmpFile.Name())
	if err != nil {
		os.Remove(tmpFile.Name())
		t.Fatalf("failed to create database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpFile.Name())
	}

	return db, cleanup
}

func TestNew(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	if db == nil {
		t.Fatal("expected database to be non-nil")
	}
}

func TestSaveMessage(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	msg := &Message{
		Sender:     "sender@example.com",
		Recipients: "recipient@example.com",
		Subject:    "Test Subject",
		Body:       "Test body content",
		RawData:    []byte("Subject: Test Subject\r\n\r\nTest body content"),
		Size:       42,
		ClientIP:   "127.0.0.1",
		IsRead:     false,
	}

	err := db.SaveMessage(msg)
	if err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	if msg.ID == 0 {
		t.Error("expected message ID to be set after save")
	}
}

func TestGetMessages(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Save a few messages
	for i := 0; i < 3; i++ {
		msg := &Message{
			Sender:     "sender@example.com",
			Recipients: "recipient@example.com",
			Subject:    "Test Subject",
			Body:       "Test body",
			Size:       10,
		}
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
	}

	messages, err := db.GetMessages()
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(messages) != 3 {
		t.Errorf("expected 3 messages, got %d", len(messages))
	}
}

func TestGetMessage(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	original := &Message{
		Sender:     "sender@example.com",
		Recipients: "recipient@example.com",
		Subject:    "Test Subject",
		Body:       "Test body content",
		RawData:    []byte("raw data"),
		Size:       100,
		ClientIP:   "192.168.1.1",
		IsRead:     false,
	}

	if err := db.SaveMessage(original); err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	retrieved, err := db.GetMessage(original.ID)
	if err != nil {
		t.Fatalf("failed to get message: %v", err)
	}

	if retrieved.Sender != original.Sender {
		t.Errorf("expected sender %q, got %q", original.Sender, retrieved.Sender)
	}
	if retrieved.Recipients != original.Recipients {
		t.Errorf("expected recipients %q, got %q", original.Recipients, retrieved.Recipients)
	}
	if retrieved.Subject != original.Subject {
		t.Errorf("expected subject %q, got %q", original.Subject, retrieved.Subject)
	}
	if retrieved.Body != original.Body {
		t.Errorf("expected body %q, got %q", original.Body, retrieved.Body)
	}
	if retrieved.Size != original.Size {
		t.Errorf("expected size %d, got %d", original.Size, retrieved.Size)
	}
	if retrieved.ClientIP != original.ClientIP {
		t.Errorf("expected client IP %q, got %q", original.ClientIP, retrieved.ClientIP)
	}
}

func TestMarkAsRead(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	msg := &Message{
		Sender:     "sender@example.com",
		Recipients: "recipient@example.com",
		Subject:    "Test",
		Body:       "Body",
		IsRead:     false,
	}

	if err := db.SaveMessage(msg); err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	if err := db.MarkAsRead(msg.ID); err != nil {
		t.Fatalf("failed to mark as read: %v", err)
	}

	retrieved, err := db.GetMessage(msg.ID)
	if err != nil {
		t.Fatalf("failed to get message: %v", err)
	}

	if !retrieved.IsRead {
		t.Error("expected message to be marked as read")
	}
}

func TestDeleteMessage(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	msg := &Message{
		Sender:     "sender@example.com",
		Recipients: "recipient@example.com",
		Subject:    "Test",
		Body:       "Body",
	}

	if err := db.SaveMessage(msg); err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	if err := db.DeleteMessage(msg.ID); err != nil {
		t.Fatalf("failed to delete message: %v", err)
	}

	_, err := db.GetMessage(msg.ID)
	if err == nil {
		t.Error("expected error when getting deleted message")
	}
}

func TestDeleteAllMessages(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Save multiple messages
	for i := 0; i < 5; i++ {
		msg := &Message{
			Sender:     "sender@example.com",
			Recipients: "recipient@example.com",
			Subject:    "Test",
			Body:       "Body",
		}
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
	}

	if err := db.DeleteAllMessages(); err != nil {
		t.Fatalf("failed to delete all messages: %v", err)
	}

	messages, err := db.GetMessages()
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(messages) != 0 {
		t.Errorf("expected 0 messages, got %d", len(messages))
	}
}

func TestGetUnreadCount(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Save 3 unread messages
	for i := 0; i < 3; i++ {
		msg := &Message{
			Sender:     "sender@example.com",
			Recipients: "recipient@example.com",
			Subject:    "Test",
			Body:       "Body",
			IsRead:     false,
		}
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
	}

	// Save 2 read messages
	for i := 0; i < 2; i++ {
		msg := &Message{
			Sender:     "sender@example.com",
			Recipients: "recipient@example.com",
			Subject:    "Test",
			Body:       "Body",
			IsRead:     false,
		}
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
		_ = db.MarkAsRead(msg.ID)
	}

	count, err := db.GetUnreadCount()
	if err != nil {
		t.Fatalf("failed to get unread count: %v", err)
	}

	if count != 3 {
		t.Errorf("expected 3 unread messages, got %d", count)
	}
}

func TestSearchMessages(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	messages := []*Message{
		{Sender: "alice@example.com", Recipients: "bob@example.com", Subject: "Hello World", Body: "Test body"},
		{Sender: "bob@example.com", Recipients: "alice@example.com", Subject: "Re: Hello World", Body: "Reply body"},
		{Sender: "charlie@example.com", Recipients: "alice@example.com", Subject: "Different Topic", Body: "Something else"},
	}

	for _, msg := range messages {
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
	}

	// Search by sender
	results, err := db.SearchMessages("alice")
	if err != nil {
		t.Fatalf("failed to search messages: %v", err)
	}
	if len(results) != 3 { // alice appears in all messages (as sender or recipient)
		t.Errorf("expected 3 results for 'alice', got %d", len(results))
	}

	// Search by subject
	results, err = db.SearchMessages("Hello")
	if err != nil {
		t.Fatalf("failed to search messages: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results for 'Hello', got %d", len(results))
	}

	// Search by body
	results, err = db.SearchMessages("Something")
	if err != nil {
		t.Fatalf("failed to search messages: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for 'Something', got %d", len(results))
	}
}

func TestMessagesOrderedByCreatedAt(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Save messages with slight delays
	subjects := []string{"First", "Second", "Third"}
	for _, subject := range subjects {
		msg := &Message{
			Sender:     "sender@example.com",
			Recipients: "recipient@example.com",
			Subject:    subject,
			Body:       "Body",
		}
		if err := db.SaveMessage(msg); err != nil {
			t.Fatalf("failed to save message: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	messages, err := db.GetMessages()
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	// Should be in reverse chronological order (newest first)
	if messages[0].Subject != "Third" {
		t.Errorf("expected first message to be 'Third', got %q", messages[0].Subject)
	}
	if messages[2].Subject != "First" {
		t.Errorf("expected last message to be 'First', got %q", messages[2].Subject)
	}
}
