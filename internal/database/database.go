package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

// New opens a SQLite database at dbPath, ensures the directory exists,
// and runs the embedded schema to create tables if they don't exist.
func New(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := ensureDir(dir); err != nil {
			return nil, fmt.Errorf("ensure db dir: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("read schema: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("exec schema: %w", err)
	}

	return db, nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
