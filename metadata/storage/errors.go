package storage

import (
	"fmt"
)

// ErrFileNotFound signals that a file does not exist when calling API.GetFile.
type ErrFileNotFound struct {
	Path Path
}

// Error returns the error message.
func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("File not found %s", e.Path)
}

// ErrFileAlreadyExists signals that a file or directory already exists.
type ErrFileAlreadyExists struct {
	Path Path
}

// Error returns the error message.
func (e ErrFileAlreadyExists) Error() string {
	return fmt.Sprintf("File already exists: %s", e.Path)
}
