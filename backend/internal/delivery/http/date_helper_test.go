package http

import (
	"testing"
	"time"
)

func TestParseDateValidFormat(t *testing.T) {
	var result time.Time
	err := parseDate("2026-03-15", &result)
	if err != nil {
		t.Fatalf("parseDate returned error for valid date: %v", err)
	}

	expected := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("parseDate result = %v, want %v", result, expected)
	}
}

func TestParseDateRejectsInvalidFormat(t *testing.T) {
	var result time.Time

	invalid := []string{
		"15-03-2026",    // DD-MM-YYYY
		"2026/03/15",    // slash separator
		"03-15-2026",    // MM-DD-YYYY
		"not-a-date",    // garbage
		"",              // empty
		"2026-13-01",    // invalid month
		"2026-02-30",    // invalid day
	}

	for _, bad := range invalid {
		err := parseDate(bad, &result)
		if err == nil {
			t.Errorf("parseDate accepted invalid date %q", bad)
		}
	}
}

func TestParseDateWritesToDestination(t *testing.T) {
	var dest time.Time
	err := parseDate("2025-12-31", &dest)
	if err != nil {
		t.Fatalf("parseDate returned error: %v", err)
	}

	if dest.Year() != 2025 || dest.Month() != 12 || dest.Day() != 31 {
		t.Errorf("dest = %v, want 2025-12-31", dest)
	}
}
