package aviation_edge

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// logUnexpectedResponse writes unexpected API responses to a debug log file
// Returns the path to the log file for inclusion in error messages
func logUnexpectedResponse(body []byte, statusCode int, endpoint string) string {
	// Create logs directory if it doesn't exist
	logDir := "/tmp/aviation_edge_errors"
	os.MkdirAll(logDir, 0755)

	// Generate unique filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("error_%s_status%d.log", timestamp, statusCode)
	logPath := filepath.Join(logDir, filename)

	// Write response with metadata
	content := fmt.Sprintf("Timestamp: %s\nStatus Code: %d\nEndpoint: %s\n\n%s\n",
		time.Now().Format(time.RFC3339),
		statusCode,
		endpoint,
		string(body))

	if err := os.WriteFile(logPath, []byte(content), 0644); err != nil {
		// Fallback if logging fails - just return temp path
		return "/tmp/aviation_edge_error.log (write failed)"
	}

	return logPath
}
