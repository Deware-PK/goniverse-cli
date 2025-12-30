package processor

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"path/filepath"
	"strings"
)

// BookMeta holds metadata information about a book.
type BookMeta struct {
	Title 	  string
	Author    string
	CoverPath string
}

type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

type Metadata struct {
	Title   []string `xml:"title"`
	Creator []string `xml:"creator"`
	Meta    []Meta   `xml:"meta"`
}

type Item struct {
	ID         string `xml:"id,attr"`
	Href       string `xml:"href,attr"`
	Properties string `xml:"properties,attr"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Package struct {
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
}


func extractMetadata(reader *zip.Reader) (*BookMeta, error) {
	var opfFile *zip.File
	for _, f := range reader.File {
		if strings.HasSuffix(f.Name, ".opf") {
			opfFile = f
			break
		}
	}

	if opfFile == nil {
		return &BookMeta{Title: "Unknown Title", Author: "Unknown Author"}, nil
	}

	rc, err := opfFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	byteValue, _ := io.ReadAll(rc)

	var pkg Package
	if err := xml.Unmarshal(byteValue, &pkg); err != nil {
		return nil, err
	}

	meta := &BookMeta{
		Title:  "Unknown Title",
		Author: "Unknown Author",
	}

	if len(pkg.Metadata.Title) > 0 {
		meta.Title = pkg.Metadata.Title[0]
	}
	if len(pkg.Metadata.Creator) > 0 {
		meta.Author = pkg.Metadata.Creator[0]
	}

	for _, item := range pkg.Manifest.Items {
		if strings.Contains(item.Properties, "cover-image") {
			meta.CoverPath = resolvePath(opfFile.Name, item.Href)
			
			return meta, nil
		}
	}

	var coverID string
	for _, m := range pkg.Metadata.Meta {
		if m.Name == "cover" {
			coverID = m.Content
			break
		}
	}
	if coverID != "" {
		for _, item := range pkg.Manifest.Items {
			if item.ID == coverID {
				meta.CoverPath = resolvePath(opfFile.Name, item.Href)
				break
			}
		}
	}

	return meta, nil
}


func resolvePath(opfPath, href string) string {
	baseDir := filepath.Dir(opfPath)
	if baseDir == "." {
		return href
	}

	return baseDir + "/" + href
}