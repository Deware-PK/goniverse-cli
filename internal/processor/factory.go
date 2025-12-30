package processor

import (
	"fmt"
	"path/filepath"
	"strings"
	"github.com/Deware-PK/goniverse-cli/internal/cleaner"
)

func NewConverter(path string, c *cleaner.HTMLCleaner) (Converter, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".epub":
		return &EpubProcessor{Cleaner: c}, nil
		
	case ".txt", ".md":
		return &TextProcessor{Cleaner: c}, nil

	default:
		return nil, fmt.Errorf("Sorry, %s is not supported.", ext)
	}
}