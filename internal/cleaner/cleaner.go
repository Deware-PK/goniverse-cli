package cleaner

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Deware-PK/goniverse-cli/internal/config"
	"golang.org/x/net/html"
)

type HTMLCleaner struct {
	Cfg *config.Config
}

func NewHTMLCleaner(cfg *config.Config) *HTMLCleaner {
	return &HTMLCleaner{Cfg: cfg}
}

func (c *HTMLCleaner) Clean(input string) string {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return input
	}

	var buf bytes.Buffer
	c.traverse(doc, &buf)

	return buf.String()
}

func (c *HTMLCleaner) traverse(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.ElementNode && c.Cfg.AllowedTags[n.Data] {
		buf.WriteString("<" + n.Data)

		for _, attr := range n.Attr {
			key := attr.Key
			val := attr.Val

			if key == "src" || key == "href" || key == "alt" || key == "title" {
				if n.Data == "img" && key == "src" {
					_, filename := filepath.Split(val)
					val = "images/" + filename
				}

				buf.WriteString(fmt.Sprintf(` %s="%s"`, key, val))
			}
		}

		buf.WriteString(">")
	} else if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.traverse(child, buf)
	}

	if n.Type == html.ElementNode && c.Cfg.AllowedTags[n.Data] {
		if n.Data != "img" && n.Data != "br" && n.Data != "hr" {
			buf.WriteString("</" + n.Data + ">")
		}
	}
}

func (c *HTMLCleaner) WrapHTML(title string, content string, currentChapter int, totalChapters int) string {
	prevLink := `<span class="disabled">&laquo; Prev</span>`
	if currentChapter > 1 {
		prevLink = fmt.Sprintf("<a href='chapter_%03d.html' class='nav-btn'>&laquo; Prev</a>", currentChapter-1)
	}

	nextLink := `<span class="disabled">Next &raquo;</span>`
	if currentChapter < totalChapters {
		nextLink = fmt.Sprintf("<a href='chapter_%03d.html' class='nav-btn'>Next &raquo;</a>", currentChapter+1)
	}

	headAndScript := `
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <link href="https://fonts.googleapis.com/css2?family=Sarabun:wght@300;400;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #ffffff;
            --text-color: #333333;
            --container-bg: #ffffff;
            --link-color: #007bff;
            --border-color: #eeeeee;
        }

        [data-theme="dark"] {
            --bg-color: #121212;
            --text-color: #e0e0e0;
            --container-bg: #1e1e1e;
            --link-color: #66b3ff;
            --border-color: #333333;
        }

        [data-theme="sepia"] {
            --bg-color: #f4ecd8;
            --text-color: #5b4636;
            --container-bg: #fdf6e3;
            --link-color: #d35400;
            --border-color: #eaddcf;
        }

        body {
            font-family: 'Sarabun', sans-serif;
            background-color: var(--bg-color);
            color: var(--text-color);
            margin: 0;
            padding: 20px;
            transition: background-color 0.3s, color 0.3s;
            font-size: 18px;
            line-height: 1.8;
        }

        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: var(--container-bg);
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.05);
        }

        img {
            max-width: 100%;
            height: auto;
            display: block;
            margin: 20px auto;
            border-radius: 4px;
        }

        /* --- Toolbar Styles --- */
        .toolbar {
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            background: var(--container-bg);
            border: 1px solid var(--border-color);
            padding: 10px 15px; /* ลด padding นิดหน่อย */
            border-radius: 50px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
            display: flex;
            gap: 10px;
            z-index: 1000;
            align-items: center;
            transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1); /* Animation นุ่มๆ */
        }

        .toolbar-content {
            display: flex;
            gap: 10px;
            align-items: center;
            opacity: 1;
            transition: opacity 0.2s ease-in-out;
        }


        .toolbar.minimized {
            bottom: auto;
            left: auto;
            transform: none;
            top: 20px;
            right: 20px;
            width: 42px;
            height: 42px;
            padding: 0;
            border-radius: 50%;
            justify-content: center;
			border-color: transparent;
			background: var(--link-color);
			color: white;
        }

        .toolbar.minimized .toolbar-content {
            display: none;
            opacity: 0;
        }
		
		.toolbar.minimized .toggle-btn {
			color: white;
			transform: rotate(180deg);
		}

        .separator {
             border-left:1px solid var(--border-color); 
             height:20px; 
             margin:0 5px;
             display: inline-block;
        }

        .btn {
            background: none;
            border: 1px solid var(--border-color);
            color: var(--text-color);
            padding: 5px 15px;
            border-radius: 20px;
            cursor: pointer;
            font-size: 14px;
            font-family: 'Sarabun', sans-serif;
            transition: all 0.2s;
        }
        
        .btn:hover { opacity: 0.8; background-color: var(--border-color); }
        .btn.active { border-color: var(--link-color); background-color: var(--link-color); color: white; }

        .toggle-btn {
            border: none;
            background: none;
            font-size: 16px;
            padding: 5px;
            display: flex;
            align-items: center;
            justify-content: center;
			color: var(--text-color);
        }
		.toggle-btn:hover { background: none; opacity: 0.7; }


        .nav-links {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid var(--border-color);
            display: flex;
            justify-content: space-between;
        }

        .nav-btn {
            padding: 10px 20px;
            background: var(--border-color);
            color: var(--text-color);
            border-radius: 5px;
            text-decoration: none;
            transition: background-color 0.2s;
        }
        .nav-btn:hover { background-color: #d0d0d0; }
        
        .disabled { opacity: 0.3; cursor: not-allowed; }

        @media (max-width: 600px) {
            body { padding: 10px; }
            .container { padding: 15px; }
            .toolbar:not(.minimized) { width: 90%; justify-content: space-between; bottom: 15px; gap: 5px; padding: 10px; }
            .toolbar-content { gap: 5px; }
            .btn { padding: 5px 10px; font-size: 12px; }
        }
    </style>

    <script>
        function setTheme(theme) {
            document.documentElement.setAttribute("data-theme", theme);
            localStorage.setItem("goniverse_theme", theme);
            document.querySelectorAll('.theme-btn').forEach(btn => btn.classList.remove('active'));
            document.getElementById('btn-' + theme).classList.add('active');
        }

        function setFontSize(size) {
            let px = '18px';
            if(size === 'small') px = '16px';
            if(size === 'large') px = '22px';
            document.body.style.fontSize = px;
            localStorage.setItem("goniverse_font", size);
        }

        function toggleToolbar() {
            const toolbar = document.getElementById('settings-toolbar');
            toolbar.classList.toggle('minimized');
            
            const isMinimized = toolbar.classList.contains('minimized');
            localStorage.setItem("goniverse_toolbar_minimized", isMinimized);
        }

        (function() {
            const savedTheme = localStorage.getItem("goniverse_theme") || 'light';
            setTheme(savedTheme);
            const savedFont = localStorage.getItem("goniverse_font") || 'medium';
            setFontSize(savedFont);

            const isMinimized = localStorage.getItem("goniverse_toolbar_minimized") === 'true';
            if (isMinimized) { toggleToolbar(); }
        })();
    </script>`

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>%s</title>
    %s
</head>
<body>

    <div class="container">
        <h1>%s</h1>
        %s
        <div class='nav-links'>
            %s
            <a href='index.html' class='nav-btn'>List of contents</a>
            %s
        </div>
    </div>

    <div class="toolbar" id="settings-toolbar">
        <button class="btn toggle-btn" onclick="toggleToolbar()" title="Toggle Settings">⚙️</button>
        
        <div class="toolbar-content">
            <button id="btn-light" class="btn theme-btn" onclick="setTheme('light')">💡</button>
            <button id="btn-sepia" class="btn theme-btn" onclick="setTheme('sepia')">📜</button>
            <button id="btn-dark" class="btn theme-btn" onclick="setTheme('dark')">🌙</button>
            <span class="separator"></span>
            <button class="btn font-btn" onclick="setFontSize('small')">A-</button>
            <button class="btn font-btn" onclick="setFontSize('large')">A+</button>
        </div>
    </div>

</body>
</html>`, title, headAndScript, title, content, prevLink, nextLink)
}
