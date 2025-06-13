// Lab 7: Implement a SQLite video metadata service

package web

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteVideoMetadataService struct{db *sql.DB}

// Uncomment the following line to ensure SQLiteVideoMetadataService implements VideoMetadataService
var _ VideoMetadataService = (*SQLiteVideoMetadataService)(nil)

// Constructor + init table
func NewSQLiteVideoMetadataService(dbPath string) (*SQLiteVideoMetadataService, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Create table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS videos (
		id TEXT PRIMARY KEY,
		uploaded_at DATETIME
	)`
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteVideoMetadataService{db: db}, nil
}

// Create inserts a new metadata row
func (s *SQLiteVideoMetadataService) Create(videoId string, uploadedAt time.Time) error {
	_, err := s.db.Exec("INSERT INTO videos (id, uploaded_at) VALUES (?, ?)", videoId, uploadedAt)
	return err
}

// List returns all video metadata
func (s *SQLiteVideoMetadataService) List() ([]VideoMetadata, error) {
	rows, err := s.db.Query("SELECT id, uploaded_at FROM videos ORDER BY uploaded_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []VideoMetadata
	for rows.Next() {
		var v VideoMetadata
		var ts string
		err := rows.Scan(&v.Id, &ts)
		if err != nil {
			return nil, err
		}
		v.UploadedAt, err = time.Parse(time.RFC3339, ts)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// Read retrieves one video metadata by ID
func (s *SQLiteVideoMetadataService) Read(videoId string) (*VideoMetadata, error) {
	row := s.db.QueryRow("SELECT id, uploaded_at FROM videos WHERE id = ?", videoId)

	var v VideoMetadata
	var ts string
	if err := row.Scan(&v.Id, &ts); err != nil {
		return nil, err
	}

	parsed, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %w", err)
	}
	v.UploadedAt = parsed
	return &v, nil
}

