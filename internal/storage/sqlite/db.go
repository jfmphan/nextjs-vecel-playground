// Package sqlite provides libSQL/SQLite-backed implementations of the domain
// repository interfaces, plus database connection and migration helpers.
package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql" // registers the "libsql" driver
	_ "modernc.org/sqlite"                                // registers the pure-Go "sqlite" driver
)

// Open connects to the database described by url. Two URL schemes are supported,
// sharing an identical SQL dialect so only the connection differs:
//
//	"file:..."                                 local SQLite via pure-Go modernc (dev)
//	"libsql://" / "http(s)://" / "ws(s)://"    Turso/libSQL remote (prod)
func Open(url, authToken string) (*sql.DB, error) {
	driverName, dsn, isFile, err := resolveDriver(url, authToken)
	if err != nil {
		return nil, err
	}
	if isFile {
		if err := ensureParentDir(dsn); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if isFile {
		// A single SQLite file serializes writes; one connection avoids
		// SQLITE_BUSY contention for this low-traffic, single-user app.
		db.SetMaxOpenConns(1)
	}
	return db, nil
}

func resolveDriver(url, authToken string) (driverName, dsn string, isFile bool, err error) {
	switch {
	case strings.HasPrefix(url, "file:"):
		return "sqlite", strings.TrimPrefix(url, "file:"), true, nil
	case hasAnyPrefix(url, "libsql://", "http://", "https://", "ws://", "wss://"):
		return "libsql", withAuthToken(url, authToken), false, nil
	default:
		return "", "", false, fmt.Errorf("unsupported DATABASE_URL scheme: %q", url)
	}
}

func withAuthToken(url, token string) string {
	if token == "" {
		return url
	}
	sep := "?"
	if strings.Contains(url, "?") {
		sep = "&"
	}
	return url + sep + "authToken=" + token
}

func ensureParentDir(dsn string) error {
	path := dsn
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

func hasAnyPrefix(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
