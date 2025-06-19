package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

type Todo struct {
	ID        int       `json:"id"`
	ChannelID string    `json:"channel_id"`
	MessageID string    `json:"message_id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type ListMessage struct {
	ID        int       `json:"id"`
	ChannelID string    `json:"channel_id"`
	MessageID string    `json:"message_id"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db := &Database{conn: conn}

	err = db.createTables()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) createTables() error {
	todoTable := `CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		channel_id TEXT NOT NULL,
		message_id TEXT NOT NULL,
		content TEXT NOT NULL,
		completed BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	listTable := `CREATE TABLE IF NOT EXISTS list_messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		channel_id TEXT NOT NULL,
		message_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.conn.Exec(todoTable)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(listTable)
	return err
}

func (db *Database) InsertTodo(channelID, messageID, content string) (int64, error) {
	query := "INSERT INTO todos (channel_id, message_id, content) VALUES (?, ?, ?)"
	result, err := db.conn.Exec(query, channelID, messageID, content)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) CompleteTodo(todoID string) (*Todo, error) {
	updateQuery := "UPDATE todos SET completed = 1 WHERE id = ?"
	_, err := db.conn.Exec(updateQuery, todoID)
	if err != nil {
		return nil, err
	}

	selectQuery := "SELECT id, channel_id, message_id, content, completed, created_at FROM todos WHERE id = ?"
	row := db.conn.QueryRow(selectQuery, todoID)

	var todo Todo
	err = row.Scan(&todo.ID, &todo.ChannelID, &todo.MessageID, &todo.Content, &todo.Completed, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (db *Database) GetTodosByChannel(channelID string) ([]Todo, error) {
	query := "SELECT id, channel_id, message_id, content, completed, created_at FROM todos WHERE channel_id = ? ORDER BY created_at DESC"
	rows, err := db.conn.Query(query, channelID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle error if needed, but typically defer will not return an error
			// since Close() is expected to succeed unless the connection is already closed.
		}
	}(rows)

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.ChannelID, &todo.MessageID, &todo.Content, &todo.Completed, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (db *Database) GetTodoMessageIDs(channelID string) ([]string, error) {
	query := "SELECT message_id FROM todos WHERE channel_id = ?"
	rows, err := db.conn.Query(query, channelID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle error if needed, but typically defer will not return an error
			// since Close() is expected to succeed unless the connection is already closed.
		}
	}(rows)

	var messageIDs []string
	for rows.Next() {
		var messageID string
		err := rows.Scan(&messageID)
		if err != nil {
			return nil, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	return messageIDs, nil
}

func (db *Database) GetListMessageIDs(channelID string) ([]string, error) {
	query := "SELECT message_id FROM list_messages WHERE channel_id = ?"
	rows, err := db.conn.Query(query, channelID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// Handle error if needed, but typically defer will not return an error
			// since Close() is expected to succeed unless the connection is already closed.
		}
	}(rows)

	var messageIDs []string
	for rows.Next() {
		var messageID string
		err := rows.Scan(&messageID)
		if err != nil {
			return nil, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	return messageIDs, nil
}

func (db *Database) DeleteTodosByChannel(channelID string) (int64, error) {
	query := "DELETE FROM todos WHERE channel_id = ?"
	result, err := db.conn.Exec(query, channelID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (db *Database) DeleteListMessagesByChannel(channelID string) error {
	query := "DELETE FROM list_messages WHERE channel_id = ?"
	_, err := db.conn.Exec(query, channelID)
	return err
}

func (db *Database) InsertListMessage(channelID, messageID string) error {
	query := "INSERT INTO list_messages (channel_id, message_id) VALUES (?, ?)"
	_, err := db.conn.Exec(query, channelID, messageID)
	return err
}

func (db *Database) Close() error {
	return db.conn.Close()
}
