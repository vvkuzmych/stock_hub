package sqlutil

import (
	"fmt"
	"strings"
)

// ConvertPlaceholders converts SQLite placeholders (?) to PostgreSQL placeholders ($1, $2, ...)
func ConvertPlaceholders(query string, driver string) string {
	if driver != "postgres" {
		return query
	}

	// Count placeholders
	count := strings.Count(query, "?")
	if count == 0 {
		return query
	}

	// Replace each ? with $1, $2, $3, etc.
	result := query
	for i := 1; i <= count; i++ {
		result = strings.Replace(result, "?", fmt.Sprintf("$%d", i), 1)
	}

	return result
}
