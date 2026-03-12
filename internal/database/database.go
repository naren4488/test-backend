package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	for _, stmt := range splitStatements(string(schema)) {
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("exec schema: %w", err)
		}
	}

	return db, nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// splitStatements returns non-empty SQL statements (split by ";"), trimmed of whitespace and comment-only lines.
func splitStatements(schema string) []string {
	var out []string
	for _, s := range strings.Split(schema, ";") {
		stmt := strings.TrimSpace(s)
		// Skip empty and comment-only blocks
		lines := strings.Split(stmt, "\n")
		var nonComment []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "--") {
				continue
			}
			nonComment = append(nonComment, line)
		}
		if len(nonComment) > 0 {
			out = append(out, strings.Join(nonComment, "\n"))
		}
	}
	return out
}
