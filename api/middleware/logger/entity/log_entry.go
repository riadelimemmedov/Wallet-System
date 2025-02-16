package entity

import "time"

//! LogEntry represents a single HTTP request log with timing and response details
type LogEntry struct {
	Timestamp time.Time
	Method    string
	Path      string
	ClientIP  string
	Status    int
	Size      int
	Latency   time.Duration
	Errors    string
}

//! NewLogEntry creates a new log entry with current timestamp
func NewLogEntry(
	method string,
	path string,
	clientIP string,
	status int,
	size int,
	latency time.Duration,
	errors string,
) LogEntry {
	return LogEntry{
		Timestamp: time.Now(),
		Method:    method,
		Path:      path,
		ClientIP:  clientIP,
		Status:    status,
		Size:      size,
		Latency:   latency,
		Errors:    errors,
	}
}
