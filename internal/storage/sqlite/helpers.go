package sqlite

import (
	"database/sql"
	"time"

	"homeinventory/internal/domain"
)

// scanner is implemented by both *sql.Row and *sql.Rows, so row-mapping helpers
// can be shared between single-row and multi-row queries.
type scanner interface {
	Scan(dest ...any) error
}

func ptrInt64(n sql.NullInt64) *int64 {
	if n.Valid {
		v := n.Int64
		return &v
	}
	return nil
}

func ptrInt(n sql.NullInt64) *int {
	if n.Valid {
		v := int(n.Int64)
		return &v
	}
	return nil
}

// dateArg converts an optional Date into a value safe to pass as a SQL argument
// (nil for NULL). It avoids a nil-pointer dereference that would occur if a nil
// *Date reached database/sql's Valuer path.
func dateArg(d *domain.Date) any {
	if d == nil || d.IsZero() {
		return nil
	}
	return *d
}

func optionalDate(d domain.Date) *domain.Date {
	if d.IsZero() {
		return nil
	}
	return &d
}

// parseTimestamp parses the formats SQLite/libSQL return for datetime('now')
// columns, falling back to the zero time if none match.
func parseTimestamp(s string) time.Time {
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
