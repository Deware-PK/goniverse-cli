
# Goniverse CLI 🪐

  

>  **Your E-book Universe, Cleaned & Delivered.** > A blazingly fast, security-focused e-book to web reader generator.

  

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)

![License](https://img.shields.io/badge/License-MIT-green.svg)

  

**Goniverse CLI** is a high-performance utility designed to transform digital books into beautiful, self-contained web experiences. Unlike traditional converters that just output raw HTML, Goniverse generates a fully functional **E-book Reader** with a modern interface, smart navigation, and security at its core.

  

Built with **Go**, it processes large libraries in seconds and ensures your reading experience is safe, fast, and accessible on any device with a browser.

  

# Features

  

### Blazingly Fast & Efficient

-  **Go Concurrency:** Utilizes advanced Goroutines and Worker Pools to process text and images simultaneously.

-  **Zero Dependencies:** Generates static HTML files that load instantly. No backend server or database required.

  

### Security First (Sanitized & Safe)

-  **Image Sanitization:** Automatically decodes and re-encodes every image to strip potential spyware, steganography, or malicious payloads.

-  **Clean Code:** Removes risky scripts and invasive tracking tags from the original e-book files.

  

### Immersive Reading Experience

-  **Real Metadata:** Extracts and displays the actual **Book Cover**, Title, and Author information.

-  **Smart Toolbar:** A collapsible floating toolbar to adjust settings without cluttering the screen.

 **Theme Engine:** 
- ☀️ **Light:** Classic paper feel.

- 🌙 **Dark:** OLED-friendly night mode.

- 📜 **Sepia:** Warm tones for long reading sessions (Novel mode).

**Customization:** Adjustable font sizes (A+/A-) and responsive layout for Mobile, Tablet, and Desktop.

  

### Future-Proof (Universal Support)

-  **Current Support:** EPUB (.epub)

-  **Coming Soon:** Kindle formats (.azw3, .mobi) and Comic archives (.cbz).

-  **PWA Ready:** Designed to be installed as a Progressive Web App (PWA) on mobile devices.

  

# Installation & Usage

### Option 1: Run from Source

  

If you have Go installed on your machine:

1. Clone the repository:

```

git clone https://github.com/Deware-PK/goniverse-cli.git

cd goniverse-cli

```

2. Install dependencies:

```

go mod tidy

```

3. Run the program:

```

go run main.go

```

  

### Option 2: Build Binary

  

You can build a standalone executable file:

```

# For Windows
go build -o goniverse-cli.exe main.go

# For Mac/Linux
go build -o goniverse-cli main.go

```


# Disclaimer

**Please Read Carefully:**

  

This tool is developed for **personal use and educational purposes only**. It is designed to format DRM-free EPUB files to improve the reading experience on web browsers.

  

-  **No DRM Removal:** This tool **does not** and **cannot** bypass or remove Digital Rights Management (DRM) protection. It only works on DRM-free files.

-  **Copyright Respect:** The author (**Deware**) does not endorse piracy or the unauthorized distribution of copyrighted material. Please ensure you own the rights to any content you process with this tool.

  

## License


This project is licensed under the MIT License - see the [LICENSE](https://github.com/Deware-PK/goniverse-cli/blob/main/LICENSE) file for details.