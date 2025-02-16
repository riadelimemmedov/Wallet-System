package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// FileStorage handles log file operations with a file handle and logger instance
type FileStorage struct {
	file   *os.File
	logger *log.Logger
}

// NewFileStorage creates a new storage instance, initializing the log file and directory
// It sets up both the file handle and logger for writing logs
func NewFileStorage(filename string) (*FileStorage, error) {
	if err := initializeLogFile(filename); err != nil {
		return nil, err
	}
	if err := ensureLogDirectory(filename); err != nil {
		return nil, err
	}
	file, err := openLogFile(filename)
	if err != nil {
		return nil, err
	}
	return &FileStorage{
		file:   file,
		logger: log.New(file, "", 0),
	}, nil
}

// Write adds a new log message to the file using the logger
// The call depth of 2 ensures correct source file reporting
func (fs *FileStorage) Write(message string) error {
	return fs.logger.Output(2, message)
}

// Close properly closes the log file handle
func (fs *FileStorage) Close() error {
	return fs.file.Close()
}

// ensureLogDirectory creates all necessary parent directories for the log file
// with appropriate permissions (0755 - rwxr-xr-x)
func ensureLogDirectory(filename string) error {
	dir := filepath.Dir(filename)
	return os.MkdirAll(dir, 0755)
}

// initializeLogFile checks if the log file exists and creates it if it doesn't
// Returns nil if file already exists
func initializeLogFile(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return createNewLogFile(filename)
	}
	return nil
}

// createNewLogFile creates a new log file with a timestamp header
// Used when initializing a new log file for the first time
func createNewLogFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer file.Close()

	header := fmt.Sprintf("=== Log file created at %s ===\n",
		time.Now().Format("2006-01-02 15:04:05"))

	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("failed to write log header: %v", err)
	}
	return nil
}

// openLogFile opens an existing log file in append and write-only mode
// with 0644 (rw-r--r--) permissions
func openLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
}
