package domain

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// Date is a timezone-free calendar date (year-month-day), used for fields like
// purchase and expiry dates where the time of day is irrelevant. It stores and
// compares as an ISO-8601 "YYYY-MM-DD" string, which also sorts chronologically
// in SQLite — so "expiring before X" is a plain string comparison.
type Date struct {
	t time.Time
}

// ParseDate parses an ISO-8601 "YYYY-MM-DD" string into a Date.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date %q, expected YYYY-MM-DD", s)
	}
	return Date{t: t}, nil
}

// Today returns the current local calendar date.
func Today() Date {
	now := time.Now()
	return Date{t: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)}
}

// AddDays returns the date shifted by n days (n may be negative).
func (d Date) AddDays(n int) Date { return Date{t: d.t.AddDate(0, 0, n)} }

// Compare returns -1, 0, or +1 as d is before, equal to, or after o.
func (d Date) Compare(o Date) int { return d.t.Compare(o.t) }

// String renders the date as "YYYY-MM-DD".
func (d Date) String() string { return d.t.Format(dateLayout) }

// IsZero reports whether the date is unset.
func (d Date) IsZero() bool { return d.t.IsZero() }

// MarshalJSON renders the date as a JSON string.
func (d Date) MarshalJSON() ([]byte, error) { return []byte(`"` + d.String() + `"`), nil }

// UnmarshalJSON parses a JSON "YYYY-MM-DD" string (empty/null becomes zero).
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*d = Date{}
		return nil
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// Value implements driver.Valuer, storing the date as an ISO string (or NULL).
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.String(), nil
}

// Scan implements sql.Scanner for reading the date back from the database.
func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*d = Date{}
		return nil
	case string:
		return d.scanString(v)
	case []byte:
		return d.scanString(string(v))
	case time.Time:
		*d = Date{t: time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Date", src)
	}
}

func (d *Date) scanString(s string) error {
	if s == "" {
		*d = Date{}
		return nil
	}
	// Some drivers return a full timestamp; keep only the date portion.
	if len(s) > len(dateLayout) {
		s = s[:len(dateLayout)]
	}
	parsed, err := ParseDate(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
