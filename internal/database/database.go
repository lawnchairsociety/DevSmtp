package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

type Message struct {
	ID         int64
	Sender     string
	Recipients string
	Subject    string
	Body       string
	RawData    []byte
	Size       int
	ClientIP   string
	IsRead     bool
	CreatedAt  time.Time
}

func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender TEXT NOT NULL,
		recipients TEXT NOT NULL,
		subject TEXT,
		body TEXT,
		raw_data BLOB,
		size INTEGER NOT NULL DEFAULT 0,
		client_ip TEXT,
		is_read BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
	CREATE INDEX IF NOT EXISTS idx_messages_is_read ON messages(is_read);
	`

	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) SaveMessage(msg *Message) error {
	query := `
	INSERT INTO messages (sender, recipients, subject, body, raw_data, size, client_ip, is_read, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query,
		msg.Sender,
		msg.Recipients,
		msg.Subject,
		msg.Body,
		msg.RawData,
		msg.Size,
		msg.ClientIP,
		msg.IsRead,
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = id

	return nil
}

func (db *DB) GetMessages() ([]Message, error) {
	query := `
	SELECT id, sender, recipients, subject, body, raw_data, size, client_ip, is_read, created_at
	FROM messages
	ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var clientIP sql.NullString
		err := rows.Scan(
			&msg.ID,
			&msg.Sender,
			&msg.Recipients,
			&msg.Subject,
			&msg.Body,
			&msg.RawData,
			&msg.Size,
			&clientIP,
			&msg.IsRead,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if clientIP.Valid {
			msg.ClientIP = clientIP.String
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (db *DB) GetMessage(id int64) (*Message, error) {
	query := `
	SELECT id, sender, recipients, subject, body, raw_data, size, client_ip, is_read, created_at
	FROM messages
	WHERE id = ?
	`

	var msg Message
	var clientIP sql.NullString
	err := db.conn.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.Sender,
		&msg.Recipients,
		&msg.Subject,
		&msg.Body,
		&msg.RawData,
		&msg.Size,
		&clientIP,
		&msg.IsRead,
		&msg.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if clientIP.Valid {
		msg.ClientIP = clientIP.String
	}

	return &msg, nil
}

func (db *DB) MarkAsRead(id int64) error {
	query := `UPDATE messages SET is_read = 1 WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) DeleteMessage(id int64) error {
	query := `DELETE FROM messages WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) DeleteAllMessages() error {
	query := `DELETE FROM messages`
	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) GetUnreadCount() (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE is_read = 0`
	var count int
	err := db.conn.QueryRow(query).Scan(&count)
	return count, err
}

func (db *DB) SearchMessages(term string) ([]Message, error) {
	query := `
	SELECT id, sender, recipients, subject, body, raw_data, size, client_ip, is_read, created_at
	FROM messages
	WHERE sender LIKE ? OR recipients LIKE ? OR subject LIKE ? OR body LIKE ?
	ORDER BY created_at DESC
	`

	searchTerm := "%" + term + "%"
	rows, err := db.conn.Query(query, searchTerm, searchTerm, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var clientIP sql.NullString
		err := rows.Scan(
			&msg.ID,
			&msg.Sender,
			&msg.Recipients,
			&msg.Subject,
			&msg.Body,
			&msg.RawData,
			&msg.Size,
			&clientIP,
			&msg.IsRead,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if clientIP.Valid {
			msg.ClientIP = clientIP.String
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
