package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/Deware-PK/goniverse-cli/internal/cleaner"
	"github.com/Deware-PK/goniverse-cli/internal/config"
	"github.com/Deware-PK/goniverse-cli/internal/processor"
)

const (
	Version = "v1.0.0"
)

func main() {

	printBanner()

	reader := bufio.NewReader(os.Stdin)
	cfg := config.NewDefaultConfig()

	fmt.Print("Enter EPUB file path (e.g. C:\\books\\test.epub): ")
	inputPath, _ := reader.ReadString('\n')
	cfg.EpubPath = cleanInput(inputPath)

	if cfg.EpubPath == "" {
		log.Fatal("Error: EPUB path cannot be empty.")
	}

	fmt.Print("Enter Output directory (e.g. C:\\books\\output): ")
	inputOut, _ := reader.ReadString('\n')
	cfg.OutputDir = cleanInput(inputOut)

	if cfg.OutputDir == "" {
		cfg.OutputDir = "clean_output"
		fmt.Println("No output folder specified. Using default: 'clean_output'")
	}

	cfg.Workers = runtime.NumCPU()

	c := cleaner.NewHTMLCleaner(cfg)
	proc, err := processor.NewConverter(cfg.EpubPath, c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nStarting Process...\nFile: %s\nOutput: %s\nWorkers: %d\n\n", cfg.EpubPath, cfg.OutputDir, cfg.Workers)

	err = proc.Process(cfg.EpubPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nDone! Check your output folder.")

	fmt.Println("Press 'Enter' to exit...")
	reader.ReadString('\n')
}

func cleanInput(text string) string {
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "\"'")
	return text
}

func printBanner() {
	art := `
!       _____             _                                 _____ _      _____ 
!      / ____|           (_)                               / ____| |    |_   _|
!     | |  __  ___  _ __  ___   _____ _ __ ___  ___ ______| |    | |      | |  
!     | | |_ |/ _ \| '_ \| \ \ / / _ \ '__/ __|/ _ \______| |    | |      | |  
!     | |__| | (_) | | | | |\ V /  __/ |  \__ \  __/      | |____| |____ _| |_ 
!      \_____|\___/|_| |_|_| \_/ \___|_|  |___/\___|       \_____|______|_____|
!
`
	fmt.Println("====================================================================================")
	fmt.Println(art)
	fmt.Println("==================================[ v1.0.0 ]========================================")
	fmt.Println()
}
