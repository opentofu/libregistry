package storage

import (
	"fmt"
	"regexp"
	"strings"
)

var pathRe = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// Path is a reference to a directory or file in the storage.
type Path string

// Validate checks if the path does not contain invalid elements.
func (p Path) Validate() error {
	if p == "" {
		return nil
	}
	parts := strings.Split(string(p), "/")
	for _, part := range parts {
		if !pathRe.MatchString(part) {
			return fmt.Errorf("invalid path: %s", p)
		}
	}
	return nil
}

// Basename returns the base path of the current path.
func (p Path) Basename() Path {
	parts := strings.Split(string(p), "/")
	if len(parts) == 1 {
		return ""
	}
	return Path(strings.Join(parts[:len(parts)-1], "/"))
}

// Filename returns the last path element.
func (p Path) Filename() string {
	parts := strings.Split(string(p), "/")
	return parts[len(parts)-1]
}
