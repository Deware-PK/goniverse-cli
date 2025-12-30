package processor

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
)

// Read image and sanitize it by re-encoding
func sanitizeAndSaveImage(srcReader io.Reader, destPath string) error {
	// Decode. To make sure it's a valid image file
	img, format, err := image.Decode(srcReader)
	if err != nil {
		return err
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Encode. Enforce to re-encode to JPEG or PNG for safety
	switch format {
	case "png":
		return png.Encode(out, img)
	default:
		// For jpg, gif, or others, convert to jpeg (Quality 85 is good)
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 85})
	}
}


func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}