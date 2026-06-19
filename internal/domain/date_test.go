package domain

import "testing"

func TestDate_ParseAndString(t *testing.T) {
	d, err := ParseDate("2026-01-15")
	if err != nil {
		t.Fatal(err)
	}
	if d.String() != "2026-01-15" {
		t.Fatalf("got %q", d.String())
	}
	if _, err := ParseDate("not-a-date"); err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestDate_Compare(t *testing.T) {
	a, _ := ParseDate("2026-01-01")
	b, _ := ParseDate("2026-12-31")
	if a.Compare(b) >= 0 {
		t.Fatal("a should sort before b")
	}
	if a.Compare(a) != 0 {
		t.Fatal("a should equal itself")
	}
}

func TestDate_JSONRoundTrip(t *testing.T) {
	d, _ := ParseDate("2026-06-19")
	b, err := d.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"2026-06-19"` {
		t.Fatalf("got %s", b)
	}
	var out Date
	if err := out.UnmarshalJSON(b); err != nil {
		t.Fatal(err)
	}
	if out.String() != d.String() {
		t.Fatalf("round trip mismatch: %s vs %s", out.String(), d.String())
	}
}
