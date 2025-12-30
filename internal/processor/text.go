package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Deware-PK/goniverse-cli/internal/cleaner"
)

type TextProcessor struct {
	Cleaner *cleaner.HTMLCleaner
}

func (p *TextProcessor) Process(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := os.MkdirAll(p.Cleaner.Cfg.OutputDir, os.ModePerm); err != nil {
		return err
	}

	var contentBuilder strings.Builder
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			contentBuilder.WriteString(fmt.Sprintf("<p>%s</p>\n", line))
		} else {
			contentBuilder.WriteString("<br>\n")
		}
	}

	fileName := filepath.Base(path)
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	finalHTML := p.Cleaner.WrapHTML(title, contentBuilder.String(), 1, 1)

	outPath := filepath.Join(p.Cleaner.Cfg.OutputDir, "index.html")
	return os.WriteFile(outPath, []byte(finalHTML), 0644)
}