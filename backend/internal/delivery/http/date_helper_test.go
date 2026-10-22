package http

import (
	"testing"
	"time"
)

func TestParseDateValid(t *testing.T) {
	var got time.Time

	if err := parseDate("2026-10-22", &got); err != nil {
		t.Fatalf("parseDate returned error for a valid date: %v", err)
	}

	want := time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("parsed = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
	if got.Hour() != 0 || got.Minute() != 0 || got.Second() != 0 {
		t.Errorf("parsed time-of-day = %02d:%02d:%02d, want 00:00:00", got.Hour(), got.Minute(), got.Second())
	}
}

func TestParseDateLeavesDestUnchangedOnError(t *testing.T) {
	// A caller that reuses a dest variable must not see a partially-written
	// value when parsing fails part-way through.
	original := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	dest := original

	if err := parseDate("2026-13-01", &dest); err == nil {
		t.Fatal("parseDate accepted month 13")
	}
	if !dest.Equal(original) {
		t.Errorf("dest = %s, want it left unchanged at %s", dest.Format(time.RFC3339), original.Format(time.RFC3339))
	}
}

func TestParseDateInvalidInputs(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{name: "wrong layout (DD-MM-YYYY)", in: "22-10-2026"},
		{name: "invalid day for month", in: "2026-02-30"},
		{name: "invalid month", in: "2026-13-01"},
		{name: "empty string", in: ""},
		{name: "non-date text", in: "bukan-tanggal"},
		{name: "date with time suffix", in: "2026-10-22T08:58:16Z"},
		{name: "single digit month", in: "2026-1-22"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var dest time.Time
			if err := parseDate(tc.in, &dest); err == nil {
				t.Errorf("parseDate(%q) = nil error, want non-nil", tc.in)
			}
		})
	}
}

func TestParseDateEpochBoundary(t *testing.T) {
	// 1970-01-01 is the zero-ish edge of the Unix epoch; make sure the layout
	// parse does not accidentally produce year 0 or a negative timestamp.
	var got time.Time
	if err := parseDate("1970-01-01", &got); err != nil {
		t.Fatalf("parseDate returned error: %v", err)
	}
	if got.Year() != 1970 || got.Month() != time.January || got.Day() != 1 {
		t.Errorf("parsed = %s, want 1970-01-01", got.Format("2006-01-02"))
	}
}
