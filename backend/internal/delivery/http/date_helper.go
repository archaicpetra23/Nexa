package http

import "time"

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string, dest *time.Time) error {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	*dest = parsed
	return nil
}
