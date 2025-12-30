package processor

import (
	"archive/zip"
	"github.com/Deware-PK/goniverse-cli/internal/cleaner"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type EpubProcessor struct {
	Cleaner *cleaner.HTMLCleaner
}

type ChapterInfo struct {
	Index    int
	FileName string
	HTMLPath string
}

func (p *EpubProcessor) Process(epubPath string) error {
	// Create output directories
	if err := os.MkdirAll(p.Cleaner.Cfg.OutputDir, os.ModePerm); err != nil {
		return err
	}

	// Create directory for images
	imgDir := filepath.Join(p.Cleaner.Cfg.OutputDir, "images")
	if err := os.MkdirAll(imgDir, os.ModePerm); err != nil {
		return err
	}

	// Open EPUB (ZIP) file
	reader, err := zip.OpenReader(epubPath)
	if err != nil {
		return fmt.Errorf("couldn't open epub: %v", err)
	}
	defer reader.Close()

	// ---------------------------------------------------------
	// Extract Metadata
	// ---------------------------------------------------------
	meta, _ := extractMetadata(&reader.Reader)
    fmt.Printf("Book: %s by %s\n", meta.Title, meta.Author)

	coverImgFilename := ""
    if meta.CoverPath != "" {
        for _, f := range reader.File {
            if f.Name == meta.CoverPath {
                rc, _ := f.Open()
                ext := filepath.Ext(f.Name)
                coverImgFilename = "cover" + ext
                destPath := filepath.Join(p.Cleaner.Cfg.OutputDir, "images", coverImgFilename)
                sanitizeAndSaveImage(rc, destPath)
                rc.Close()
                fmt.Println("Cover image extracted.")
                break
            }
        }
    }

	// ---------------------------------------------------------
	// Separate HTML and Image files
	// ---------------------------------------------------------
	var htmlFiles []*zip.File
	var imageFiles []*zip.File

	for _, f := range reader.File {
		if isHTML(f.Name) {
			htmlFiles = append(htmlFiles, f)
		} else if isImage(f.Name) {
			imageFiles = append(imageFiles, f)
		}
	}

	sort.Slice(htmlFiles, func(i, j int) bool {
		return htmlFiles[i].Name < htmlFiles[j].Name
	})

	totalFiles := len(htmlFiles)
	fmt.Printf("Found %d HTML chapters and %d Images.\n", totalFiles, len(imageFiles))

	// ---------------------------------------------------------
	// Image Processing
	// ---------------------------------------------------------
	if len(imageFiles) > 0 {
		fmt.Println("Processing Images...")
		var imgWg sync.WaitGroup
		
		imgJobs := make(chan *zip.File, len(imageFiles))

		imgWorkers := 5
		for w := 0; w < imgWorkers; w++ {
			imgWg.Add(1)
			go func() {
				defer imgWg.Done()
				for f := range imgJobs {
					// Open file from ZIP
					rc, err := f.Open()
					if err != nil {
						continue
					}
					
					// Prepare destination path
					_, fname := filepath.Split(f.Name)
					destPath := filepath.Join(imgDir, fname)

					sanitizeAndSaveImage(rc, destPath) 
					
					rc.Close()
				}
			}()
		}

		for _, f := range imageFiles {
			imgJobs <- f
		}
		close(imgJobs)
		
		imgWg.Wait() 
		fmt.Println("Images Processed.")
	}

	// ---------------------------------------------------------
	// HTML Processing
	// ---------------------------------------------------------
	fmt.Println("Processing Chapters...")
	
	jobs := make(chan struct {
		f *zip.File
		i int
	}, totalFiles)
	
	results := make(chan ChapterInfo, totalFiles)
	
	var wg sync.WaitGroup

	// Worker Pool
	workers := p.Cleaner.Cfg.Workers
	if workers <= 0 { workers = 1 }

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// Read
				content, _ := readFileContent(job.f)
				
				// Clean
				cleanContent := p.Cleaner.Clean(content)

				// Insert
				title := fmt.Sprintf("Chapter %d", job.i+1)
				finalHTML := p.Cleaner.WrapHTML(title, cleanContent, job.i+1, totalFiles)

				// Rename index file
				outName := fmt.Sprintf("chapter_%03d.html", job.i+1)
				outPath := filepath.Join(p.Cleaner.Cfg.OutputDir, outName)

				// Save
				err := os.WriteFile(outPath, []byte(finalHTML), 0644)
				if err == nil {
					// fmt.Printf("[Worker %d] Saved: %s\n", workerID, outName) // Debug log
					results <- ChapterInfo{Index: job.i + 1, FileName: job.f.Name, HTMLPath: outName}
				}
			}
		}(w)
	}

	// Sending to workers
	for i, f := range htmlFiles {
		jobs <- struct {
			f *zip.File
			i int
		}{f, i}
	}
	close(jobs)

	// Await html workers completion
	go func() {
		wg.Wait()
		close(results)
	}()

	// ---------------------------------------------------------
	// Create Index page
	// ---------------------------------------------------------
	var chapters []ChapterInfo
	for res := range results {
		chapters = append(chapters, res)
	}

	// Sort chapters
	sort.Slice(chapters, func(i, j int) bool {
		return chapters[i].Index < chapters[j].Index
	})

	return p.generateIndex(chapters, meta, coverImgFilename)
}

func (p *EpubProcessor) generateIndex(chapters []ChapterInfo, meta *BookMeta, coverFilename string) error {
	var listItems string
	for _, ch := range chapters {
		listItems += fmt.Sprintf("<li><a href='%s'>Chapter %d : %s</a></li>", ch.HTMLPath, ch.Index, ch.FileName)
	}

	headerHTML := "<div style='text-align:center; margin-bottom:40px;'>"
	
	if coverFilename != "" {
		headerHTML += fmt.Sprintf("<img src='images/%s' style='max-width:300px; box-shadow:0 10px 20px rgba(0,0,0,0.2); border-radius:8px; margin-bottom:20px;'><br>", coverFilename)
	}
	
	headerHTML += fmt.Sprintf("<h1 style='margin-bottom:5px;'>%s</h1>", meta.Title)
	headerHTML += fmt.Sprintf("<p style='color:var(--quote-color); font-size:1.2em;'>By %s</p>", meta.Author)
	headerHTML += "</div>"

	content := fmt.Sprintf("%s<hr style='border:0; border-top:1px solid var(--border-color); margin:30px 0;'><h3>สารบัญ</h3><ul>%s</ul>", headerHTML, listItems)
	
	indexHTML := p.Cleaner.WrapHTML(meta.Title, content, -1, -1)

	return os.WriteFile(filepath.Join(p.Cleaner.Cfg.OutputDir, "index.html"), []byte(indexHTML), 0644)
}

func isHTML(name string) bool {
	return len(name) > 5 && (name[len(name)-5:] == ".html" || name[len(name)-6:] == ".xhtml")
}

func readFileContent(f *zip.File) (string, error) {
	rc, err := f.Open()
	if err != nil { return "", err }
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil { return "", err }
	return string(content), nil
}